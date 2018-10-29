package slick

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
)

type botReply struct {
	To   string
	Text string
}

// Message represents a specific slack message
type Message struct {
	*slack.Msg
	SubMessage  *slack.Msg
	bot         *Bot
	MentionsMe  bool
	IsEdit      bool
	FromMe      bool
	FromUser    *slack.User
	FromChannel *Channel

	// Match contains the result of
	// Listener.Matches.FindStringSubmatch(msg.Text), when `Matches`
	// is set on the `Listener`.
	Match []string
}

// IsPrivate determines if a message is private or not
func (msg *Message) IsPrivate() bool {
	return strings.HasPrefix(msg.Channel, "D")
}

// ContainsAnyCased searches for at least one case-sensitive word
func (msg *Message) ContainsAnyCased(strs []string) bool {
	for _, s := range strs {
		if strings.Contains(msg.Text, s) {
			return true
		}
	}
	return false
}

// ContainsAny searches for at least one noncase-sensitive matching string
func (msg *Message) ContainsAny(strs []string) bool {
	lowerStr := strings.ToLower(msg.Text)

	for _, s := range strs {
		lowerInput := strings.ToLower(s)

		if strings.Contains(lowerStr, lowerInput) {
			return true
		}
	}
	return false
}

// ContainsAll searches for all strings in a noncase-sensitive fashion
func (msg *Message) ContainsAll(strs []string) bool {

	lowerStr := strings.ToLower(msg.Text)

	for _, s := range strs {
		lowerInput := strings.ToLower(s)

		if !strings.Contains(lowerStr, lowerInput) {
			return false
		}
	}
	return true
}

// Contains searches for a single string in a noncase-sensitive fashion
func (msg *Message) Contains(s string) bool {
	lowerStr := strings.ToLower(msg.Text)
	lowerInput := strings.ToLower(s)

	if strings.Contains(lowerStr, lowerInput) {
		return true
	}
	return false
}

// HasPrefix returns true if a message starts with a given string
func (msg *Message) HasPrefix(prefix string) bool {
	return strings.HasPrefix(msg.Text, prefix)
}

// AddReaction adds a reaction to a message
func (msg *Message) AddReaction(emoticon string) *Message {
	msg.bot.Slack.AddReaction(emoticon, slack.NewRefToMessage(msg.Channel, msg.Timestamp))
	return msg
}

// RemoveReaction removes a reaction from a message
func (msg *Message) RemoveReaction(emoticon string) *Message {
	msg.bot.Slack.RemoveReaction(emoticon, slack.NewRefToMessage(msg.Channel, msg.Timestamp))
	return msg
}

// ListenReaction listens for a reaction on a message
func (msg *Message) ListenReaction(reactListen *ReactionListener) {
	msg.bot.ListenReaction(msg.Timestamp, reactListen)
}

// Reply sends a message back to the source it came from, without a mention
func (msg *Message) Reply(text string, v ...interface{}) *Reply {
	to := msg.User
	if msg.Channel != "" {
		to = msg.Channel
	}
	text = Format(text, v...)
	return msg.bot.SendOutgoingMessage(text, to)
}

// ReplyPrivately replies to the user in an IM
func (msg *Message) ReplyPrivately(text string, v ...interface{}) *Reply {
	text = Format(text, v...)
	return msg.bot.SendPrivateMessage(msg.User, text)
}

// ReplyMention replies with a @mention named prefixed, when replying
// in public. When replying in private, nothing is added.
func (msg *Message) ReplyMention(text string, v ...interface{}) *Reply {
	if msg.IsPrivate() {
		return msg.Reply(text, v...)
	}
	prefix := ""
	if msg.FromUser != nil {
		prefix = fmt.Sprintf("<@%s> ", msg.FromUser.Name)
	}
	return msg.Reply(fmt.Sprintf("%s%s", prefix, text), v...)
}

// String returns a message with field:value as a string
func (msg *Message) String() string {
	return fmt.Sprintf("%#v", msg)
}

func (msg *Message) applyMentionsMe(bot *Bot) {
	if msg.IsPrivate() {
		msg.MentionsMe = true
	}

	m := reAtMention.FindStringSubmatch(msg.Text)
	if m != nil && m[1] == bot.Myself.ID {
		msg.MentionsMe = true
	}
}

func (msg *Message) applyFromMe(bot *Bot) {
	if msg.User != "" && msg.User == bot.Myself.ID {
		msg.FromMe = true
	}
}

var reAtMention = regexp.MustCompile(`<@([A-Z0-9]+)(|([^>]+))>`)

// Format conditionally formats using fmt.Sprintf if there is more
// than one argument, otherwise returns the first parameter
// uninterpreted.
func Format(s string, v ...interface{}) string {
	count := len(v)
	if count == 0 {
		return s
	}
	return fmt.Sprintf(s, v...)
}
