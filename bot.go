// Package slick is a ChatOps framework for Slack
package slick

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cskr/pubsub"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Bot is the main slick bot instance. It is passed throughout, and
// has references to most useful objects.
type Bot struct {
	// Global bot configuration
	configFile string
	Config     SlackConfig

	// Logging configuration
	Logging Logging

	// Slack connectivity
	Slack             *slack.Client
	rtm               *slack.RTM
	Users             map[string]slack.User
	Channels          map[string]Channel
	channelUpdateLock sync.Mutex
	Myself            slack.UserDetails

	// Internal handling
	listeners     []*Listener
	addListenerCh chan *Listener
	delListenerCh chan *Listener
	outgoingMsgCh chan *slack.OutgoingMessage

	// Storage
	DB *bolt.DB

	// Inter-plugins communications. Use topics like
	// "pluginName:eventType[:someOtherThing]"
	PubSub *pubsub.PubSub

	// Other features
	WebServer WebServer
	Mood      Mood
}

// New returns a new bot instance, initialized with the provided config
// file. If an empty string is provided as the config file path, Slick
// searches the working directy and $HOME/.slick/ for a file called
// config.json|toml|yaml instead
func New(configFile string) *Bot {
	bot := &Bot{
		configFile:    configFile,
		outgoingMsgCh: make(chan *slack.OutgoingMessage, 500),
		addListenerCh: make(chan *Listener, 500),
		delListenerCh: make(chan *Listener, 500),

		Users:    make(map[string]slack.User),
		Channels: make(map[string]Channel),

		PubSub: pubsub.New(500),
	}

	http.DefaultClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	return bot
}

// Run starts the bot
func (bot *Bot) Run() {
	// Config for Slack and logging are read in
	bot.loadBaseConfig()

	// Configure logging
	err := bot.setupLogging()
	if err != nil {
		log.Fatal("Error setting up logging.")
	}

	// Write PID
	err = bot.writePID()
	if err != nil {
		log.Fatal("Couldn't write PID file:", err)
	}

	db, err := bolt.Open(bot.Config.DBPath, 0600, nil)
	if err != nil {
		log.WithError(err).Fatalf("Could not initialize BoltDB key/value store: %s", err)
	}
	defer func() {
		log.Fatal("Database is closing")
		db.Close()
	}()
	bot.DB = db

	// Init all plugins
	var enabledPlugins []string
	for _, plugin := range registeredPlugins {
		pluginType := reflect.TypeOf(plugin)
		if pluginType.Kind() == reflect.Ptr {
			pluginType = pluginType.Elem()
		}
		var typeList []string
		if _, ok := plugin.(PluginInitializer); ok {
			typeList = append(typeList, "Plugin")
		}
		if _, ok := plugin.(WebServer); ok {
			typeList = append(typeList, "WebServer")
		}
		if _, ok := plugin.(WebServerAuth); ok {
			typeList = append(typeList, "WebServerAuth")
		}
		if _, ok := plugin.(WebPlugin); ok {
			typeList = append(typeList, "WebPlugin")
		}

		log.Printf("Plugin %s implements %s", pluginType.String(),
			strings.Join(typeList, ", "))
		enabledPlugins = append(enabledPlugins, strings.Replace(pluginType.String(), ".", "_", -1))
	}

	initWebServer(bot, enabledPlugins)
	initWebPlugins(bot)

	if bot.WebServer != nil {
		go bot.WebServer.RunServer()
	}

	initChatPlugins(bot)

	bot.Slack = slack.New(bot.Config.ApiToken)
	bot.Slack.SetDebug(bot.Config.Debug)

	rtm := bot.Slack.NewRTM()
	bot.rtm = rtm

	bot.setupHandlers()

	bot.rtm.ManageConnection()
}

func (bot *Bot) writePID() error {
	var serverConf struct {
		Server struct {
			Pidfile string `mapstructure:"pid_file"`
		}
	}

	err := bot.LoadConfig(&serverConf)
	if err != nil {
		return err
	}

	if serverConf.Server.Pidfile == "" {
		return nil
	}

	pid := os.Getpid()
	pidb := []byte(strconv.Itoa(pid))
	return ioutil.WriteFile(serverConf.Server.Pidfile, pidb, 0755)
}

// Listen registers a listener for messages and events. There are two main
// handling functions on a Listener: MessageHandlerFunc and EventHandlerFunc.
// MessageHandlerFunc is filtered by a bunch of other properties of the Listener,
// whereas EventHandlerFunc will receive all events unfiltered, but with
// *slick.Message instead of a raw *slack.MessageEvent (it's in there anyway),
// which adds a bunch of useful methods to it.
//
// Explore the Listener for more details.
func (bot *Bot) Listen(listen *Listener) error {
	listen.Bot = bot

	err := listen.checkParams()
	if err != nil {
		log.Println("Bot.Listen(): Invalid Listener: ", err)
		return err
	}

	bot.addListener(listen)

	return nil
}

// ListenReaction will dispatch the listener with matching incoming reactions.
// `item` can be a timestamp or a file ID.
func (bot *Bot) ListenReaction(item string, reactListen *ReactionListener) {
	listen := reactListen.newListener()
	listen.EventHandlerFunc = func(_ *Listener, event interface{}) {
		re := ParseReactionEvent(event)
		if re == nil {
			return
		}

		if item != re.Item.Timestamp && item != re.Item.File {
			return
		}

		if re.User == bot.Myself.ID {
			return
		}

		if !reactListen.filterReaction(re) {
			return
		}

		re.Listener = reactListen

		reactListen.HandlerFunc(reactListen, re)
	}
	bot.Listen(listen)
}

func (bot *Bot) addListener(listen *Listener) {
	listen.setupChannels()
	if listen.isManaged() {
		go listen.launchManager()
	}
	bot.addListenerCh <- listen
}

func (bot *Bot) Notify(room, color, format, msg string, notify bool) error {
	log.Println("DEPRECATED. Please use the Slack API with .PostMessage")
	// bot.api.PostMessage(room, msg, slack.PostMessageParameters{
	// 	Attachments: []slack.Attachment{
	// 		{
	// 			Color: color,
	// 			Text: msg,
	// 		},
	// 	},
	// })
	return nil
}

func (bot *Bot) setupHandlers() {
	go bot.replyHandler()
	go bot.messageHandler()
	log.Println("Bot ready")
}

func (bot *Bot) cacheUsers(users []slack.User) {
	bot.Users = make(map[string]slack.User)
	for _, user := range users {
		bot.Users[user.ID] = user
	}
}

func (bot *Bot) cacheChannels(channels []slack.Channel, groups []slack.Group, ims []slack.IM) {
	log.Debugf("Channels: %v", len(channels))
	log.Debugf("Groups: %v", len(groups))
	log.Debugf("DM's: %v", len(ims))
	bot.Channels = make(map[string]Channel)
	for _, channel := range channels {
		bot.updateChannel(ChannelFromSlackChannel(channel))
	}

	for _, group := range groups {
		bot.updateChannel(ChannelFromSlackGroup(group))
	}

	for _, im := range ims {
		bot.updateChannel(ChannelFromSlackIM(im))
	}
}

func (bot *Bot) loadBaseConfig() {

	bot.readInConfig() // Find and parse the config file

	var config struct {
		Logging Logging
		Slack   SlackConfig
	}

	err := bot.LoadConfig(&config)
	if err != nil {
		log.WithError(err).Fatalln("Error loading config file.")
	}

	bot.Config = config.Slack
	bot.Logging = config.Logging
}

// readInConfig reads the config file, unmarshals the given format (JSON, YAML or TOML)
// and makes the result available for later use in LoadConfig()
func (bot *Bot) readInConfig() {

	// Use viper to find a default config file, or open the provided file is set
	if bot.configFile == "" {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")            // Look for config in the working directory
		viper.AddConfigPath("$HOME/.slick") // Look for config in .slick folder in home directory
	} else {
		// Ensure the config file cannot be read of write by others
		if err := checkPermission(bot.configFile); err != nil {
			log.WithError(err).Fatal("Error checking permissions.")
		}

		viper.SetConfigFile(bot.configFile)
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Return an error if the config file cannot be parsed
		log.WithError(err).Fatalf("fatal error config file: %s", err)
	}
}

// LoadConfig populates a given struct with the values from the config file
func (bot *Bot) LoadConfig(config interface{}) (err error) {

	err = viper.Unmarshal(&config)
	if err != nil {
		log.WithError(err).Errorf("unable to decode into struct, %v", err)
		return err
	}

	return nil
}

func (bot *Bot) replyHandler() {
	for {
		outMsg := <-bot.outgoingMsgCh
		if outMsg == nil {
			continue
		}

		bot.rtm.SendMessage(outMsg)

		time.Sleep(50 * time.Millisecond)
	}
}

// SendToChannel sends a message to a given channel
func (bot *Bot) SendToChannel(channelName string, message string) *Reply {
	channel := bot.GetChannelByName(channelName)

	if channel == nil {
		log.WithFields(log.Fields{
			"Type":    "ChannelNotFound",
			"Channel": channelName,
		}).Error("Error sending message to channel.")

		return nil
	}

	log.WithFields(log.Fields{
		"Type":    "SendingMessage",
		"Channel": channelName,
		"Message": message,
	}).Debug("Sending message to channel.")

	return bot.SendOutgoingMessage(message, channel.ID)
}

// SendOutgoingMessage schedules the message for departure and returns
// a Reply which can be listened on. See type `Reply`.
func (bot *Bot) SendOutgoingMessage(text string, to string) *Reply {
	log.WithFields(log.Fields{
		"Type":      "SendingMessage",
		"Recipient": to,
		"Message":   text,
	}).Debug("Sending outgoing message.")

	outMsg := bot.rtm.NewOutgoingMessage(text, to)
	bot.outgoingMsgCh <- outMsg

	return &Reply{outMsg, bot}
}

// SendPrivateMessage sends a message to a user
func (bot *Bot) SendPrivateMessage(username, message string) *Reply {
	user := bot.GetUser(username)
	if user == nil {
		log.WithFields(log.Fields{
			"Type":      "UserDoesNotExist",
			"Recipient": username,
			"Message":   message,
		}).Warn("Error sending message.")

		return nil
	}

	imChannel := bot.OpenIMChannelWith(user)
	if imChannel == nil {
		log.WithFields(log.Fields{
			"Type":         "IMChannelDoesNotExist",
			"Recipient":    user.Name,
			"Recipient ID": user.ID,
			"Message":      message,
		}).Warn("Error sending message.")

		return nil
	}

	log.WithFields(log.Fields{
		"Type":       "SendingPrivateMessage",
		"IM Channel": imChannel.ID,
		"Message":    message,
	}).Info("Sending private message.")

	outMsg := bot.rtm.NewOutgoingMessage(message, imChannel.ID)
	bot.outgoingMsgCh <- outMsg

	return &Reply{outMsg, bot}
}

func (bot *Bot) removeListener(listen *Listener) {
	for i, element := range bot.listeners {
		if element == listen {
			// following: https://code.google.com/p/go-wiki/wiki/SliceTricks
			copy(bot.listeners[i:], bot.listeners[i+1:])
			bot.listeners[len(bot.listeners)-1] = nil
			bot.listeners = bot.listeners[:len(bot.listeners)-1]
			return
		}
	}
}

func (bot *Bot) messageHandler() {
	for {
	nextMessages:
		select {
		case listen := <-bot.addListenerCh:
			bot.listeners = append(bot.listeners, listen)

		case listen := <-bot.delListenerCh:
			bot.removeListener(listen)

		case event := <-bot.rtm.IncomingEvents:
			bot.handleRTMEvent(&event)
		}

		// Always flush listeners deletions between messages, so a
		// Close()'d Listener never processes another message.
		for {
			select {
			case listen := <-bot.delListenerCh:
				bot.removeListener(listen)
			default:
				goto nextMessages
			}
		}
	}
}

func (bot *Bot) handleRTMEvent(event *slack.RTMEvent) {
	var msg *Message
	var client = bot.Slack
	//var reaction interface{}

	switch ev := event.Data.(type) {
	/**
	 * Connection handling...
	 */
	case *slack.LatencyReport:
		log.WithFields(log.Fields{
			"Type":    "LatencyReport",
			"Latency": ev.Value,
		}).Debug("Latency Report.")
	case *slack.RTMError:
		log.WithFields(log.Fields{
			"Type":      "RTMError",
			"ErrorCode": ev.Code,
			"Message":   ev.Msg,
		}).Error("Real Time Messenger Error.")
	case *slack.ConnectedEvent:
		// Replacing ev.Info.Channels
		channels, err := client.GetChannels(false)
		if err != nil {
			panic("SLACK DEPRECATED ANOTHER API LOL")
		}

		// Replacing ev.Info.Groups
		groups, err := client.GetGroups(false)
		if err != nil {
			panic("SLACK DEPRECATED ANOTHER API LOL")
		}

		// Replacing ev.Info.IMs
		ims, err := client.GetIMChannels()
		if err != nil {
			panic("SLACK DEPRECATED ANOTHER API LOL")
		}

		// Replacing ev.Info.Users
		users, err := client.GetUsers()
		if err != nil {
			panic("SLACK DEPRECATED ANOTHER API LOL")
		}

		log.Printf("Bot connected, connection_count=%d", ev.ConnectionCount)
		bot.Myself = *ev.Info.User
		bot.cacheUsers(users)                    // Info.Users is deprecated
		bot.cacheChannels(channels, groups, ims) // Info.Channels is deprecated

		for _, channelName := range bot.Config.JoinChannels {
			channel := bot.GetChannelByName(channelName)
			if channel != nil && !channel.IsMember {
				bot.Slack.JoinChannel(channel.ID)
			}
		}

	case *slack.DisconnectedEvent:
		log.Println("Bot disconnected")

	case *slack.ConnectingEvent:
		log.Printf("Bot connecting, connection_count=%d, attempt=%d", ev.ConnectionCount, ev.Attempt)

	case *slack.HelloEvent:
		log.Println("Got a HELLO from websocket")

	/**
	 * Message dispatch and handling
	 */
	case *slack.MessageEvent:
		log.WithFields(log.Fields{
			"Message": ev,
		}).Debug("Message received.")

		msg = &Message{
			Msg:        &ev.Msg,
			SubMessage: ev.SubMessage,
			bot:        bot,
		}

		userID := ev.User
		switch ev.Msg.SubType {
		case "message_changed":
			userID = ev.SubMessage.User
			msg.Msg.Text = ev.SubMessage.Text
			msg.IsEdit = true
		case "channel_topic":
			if channel, ok := bot.Channels[ev.Channel]; ok {
				channel.Topic = slack.Topic{
					Value:   ev.Topic,
					Creator: ev.User,
					LastSet: unixFromTimestamp(ev.Timestamp),
				}
				bot.Channels[ev.Channel] = channel
			}
		case "channel_purpose":
			if channel, ok := bot.Channels[ev.Channel]; ok {
				channel.Purpose = slack.Purpose{
					Value:   ev.Purpose,
					Creator: ev.User,
					LastSet: unixFromTimestamp(ev.Timestamp),
				}
				bot.Channels[ev.Channel] = channel
			}
		}

		// We do some heavy logging here because this is troublesome
		// to find when it breaks and Slack breaks it by deprecating API's

		user, ok := bot.Users[userID]
		if ok {
			log.Debug("User map is ok.")
			msg.FromUser = &user
		} else {
			log.WithFields(log.Fields{
				"Type":  "BrokenUserMap",
				"Users": len(bot.Users),
			}).Error("User map is broken.")
		}

		channel, ok := bot.Channels[ev.Channel]
		if ok {
			log.Debug("Channel map is ok.")
			msg.FromChannel = &channel
		} else {
			log.WithFields(log.Fields{
				"Type":     "BrokenChannelMap",
				"Channels": len(bot.Channels),
			}).Error("Channel map is broken.")
		}

		msg.applyMentionsMe(bot)
		msg.applyFromMe(bot)

	case *slack.PresenceChangeEvent:
		user := bot.Users[ev.User]
		log.Printf("User %q is now %q", user.Name, ev.Presence)
		user.Presence = ev.Presence

	// TODO: manage im_open, im_close, and im_created ?

	// case *slack.ReactionAddedEvent:
	// 	reaction = ev
	// case *slack.ReactionRemovedEvent:
	// 	reaction = ev

	/**
	 * User changes
	 */
	case *slack.UserChangeEvent:
		bot.Users[ev.User.ID] = ev.User

	/**
	 * Handle slack Channel changes
	 */
	case *slack.ChannelRenameEvent:
		channel := bot.Channels[ev.Channel.ID]
		channel.Name = ev.Channel.Name
		bot.updateChannel(channel)

	case *slack.ChannelJoinedEvent:
		bot.updateChannel(ChannelFromSlackChannel(ev.Channel))

	case *slack.ChannelCreatedEvent:
		c := Channel{}
		c.ID = ev.Channel.ID
		c.Name = ev.Channel.Name
		c.Creator = ev.Channel.Creator
		c.IsChannel = true
		bot.updateChannel(c)

		// NICE TODO: poll the API to get a full Channel object ? many
		// things are missing here

	case *slack.ChannelDeletedEvent:
		bot.deleteChannel(ev.Channel)

	case *slack.ChannelArchiveEvent:
		channel := bot.Channels[ev.Channel]
		channel.IsArchived = true
		bot.updateChannel(channel)

	case *slack.ChannelUnarchiveEvent:
		channel := bot.Channels[ev.Channel]
		channel.IsArchived = false
		bot.updateChannel(channel)

	/**
	 * Handle slack Group changes
	 */
	case *slack.GroupRenameEvent:
		group := bot.Channels[ev.Group.ID]
		group.Name = ev.Group.Name
		bot.updateChannel(group)

	case *slack.GroupJoinedEvent:
		bot.updateChannel(ChannelFromSlackChannel(ev.Channel))

	case *slack.GroupCreatedEvent:
		c := Channel{}
		c.ID = ev.Channel.ID
		c.Name = ev.Channel.Name
		c.Creator = ev.Channel.Creator
		c.IsGroup = true
		bot.updateChannel(c)

		// NICE: poll the API to get a full Group object ? many
		// things are missing here

	case *slack.GroupCloseEvent:
		bot.deleteChannel(ev.Channel)

	case *slack.GroupArchiveEvent:
		group := bot.Channels[ev.Channel]
		group.IsArchived = true
		bot.updateChannel(group)

	case *slack.GroupUnarchiveEvent:
		group := bot.Channels[ev.Channel]
		group.IsArchived = false
		bot.updateChannel(group)

	/**
	 * Handle slack IM changes
	 */
	case *slack.IMCreatedEvent:
		c := Channel{}
		c.ID = ev.Channel.ID
		c.User = ev.User
		c.IsIM = true
		bot.updateChannel(c)

	case *slack.IMOpenEvent:
		c := Channel{}
		c.ID = ev.Channel
		c.User = ev.User
		c.IsIM = true
		bot.updateChannel(c)

	case *slack.IMCloseEvent:
		bot.deleteChannel(ev.Channel)

	/**
	 * Errors
	 */
	case *slack.AckErrorEvent:
		jsonCnt, _ := json.MarshalIndent(ev, "", "  ")
		fmt.Printf("Error: %s", jsonCnt)

	default:
		log.Printf("Event: %T", ev)
		//log.Printf("Unexpected: %#v", ev)
	}

	// Dispatch listeners
	for _, listen := range bot.listeners {
		if msg != nil && listen.MessageHandlerFunc != nil {
			listen.filterAndDispatchMessage(msg)
		}

		if listen.EventHandlerFunc != nil {
			var handleEvent interface{} = event.Data
			if msg != nil {
				handleEvent = msg
			}
			listen.EventHandlerFunc(listen, handleEvent)
		}
	}

}

// Disconnect the websocket.
func (bot *Bot) Disconnect() {
	// FIXME: implement a Reconnect() method.. calling the RTM method of the same name.
	// QUERYME: do we need that, really ?
	bot.rtm.Disconnect()
}

// GetUser returns a *slack.User by ID, Name, RealName or Email
func (bot *Bot) GetUser(find string) *slack.User {
	for _, user := range bot.Users {
		//log.Printf("Hmmmm, %#v", user)
		if user.Profile.Email == find || user.ID == find || user.Name == find || user.RealName == find {
			return &user
		}
	}
	return nil
}

// GetChannelByName returns a *slack.Channel by Name
func (bot *Bot) GetChannelByName(name string) *Channel {
	name = strings.TrimLeft(name, "#")
	for _, channel := range bot.Channels {
		if channel.Name == name {
			return &channel
		}
	}
	return nil
}

// GetIMChannelWith returns the channel used to communicate with the
// specified slack user
func (bot *Bot) GetIMChannelWith(user *slack.User) *Channel {
	for _, channel := range bot.Channels {
		if !channel.IsIM {
			continue
		}
		if channel.User == user.ID {
			return &channel
		}
	}
	return nil
}

// OpenIMChannelWith opens a conversation with the given slack User
func (bot *Bot) OpenIMChannelWith(user *slack.User) *Channel {
	dmChannel := bot.GetIMChannelWith(user)
	if dmChannel != nil {
		return dmChannel
	}

	log.Printf("Opening a new IM conversation with %q (%s)", user.ID, user.Name)
	_, _, chanID, err := bot.Slack.OpenIMChannel(user.ID)
	if err != nil {
		return nil
	}

	c := Channel{
		ID:   chanID,
		IsIM: true,
		User: user.ID,
	}
	bot.updateChannel(c)

	return &c
}

func (bot *Bot) updateChannel(channel Channel) {
	bot.channelUpdateLock.Lock()
	bot.Channels[channel.ID] = channel
	bot.channelUpdateLock.Unlock()
}

func (bot *Bot) deleteChannel(id string) {
	bot.channelUpdateLock.Lock()
	delete(bot.Channels, id)
	bot.channelUpdateLock.Unlock()
}
