// Package healthy is a Slick plugin that evaluates whether URLs return 200's or not
package healthy

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/CapstoneLabs/slick"
	log "github.com/sirupsen/logrus"
)

// Healthy is a struct holding URL's to evaluate
type Healthy struct {
	urls []string
}

func init() {
	slick.RegisterPlugin(&Healthy{})
}

// InitPlugin loads configuration and listens for new messages
func (healthy *Healthy) InitPlugin(bot *slick.Bot) {
	var conf struct {
		HealthCheck struct {
			Urls []string
		}
	}

	bot.LoadConfig(&conf)

	healthy.urls = conf.HealthCheck.Urls

	bot.Listen(&slick.Listener{
		MentionsMeOnly:     true,
		ContainsAny:        []string{"health", "healthy?", "health_check"},
		MessageHandlerFunc: healthy.ChatHandler,
	})
}

// ChatHandler replies to the end user
func (healthy *Healthy) ChatHandler(listen *slick.Listener, msg *slick.Message) {
	log.Println("Health check. Requested by", msg.FromUser.Name)
	msg.Reply(healthy.CheckAll())
}

// CheckAll checks each URL in the struct
func (healthy *Healthy) CheckAll() string {
	result := make(map[string]bool)
	failed := make([]string, 0)
	for _, url := range healthy.urls {
		ok := check(url)
		result[url] = ok
		if !ok {
			failed = append(failed, url)
		}
	}
	if len(failed) == 0 {
		return "All green (For " +
			strings.Join(healthy.urls, ", ") + ")"
	} else {
		return "WARN!! Something wrong with " +
			strings.Join(failed, ", ")
	}
}

func check(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		return false
	}
	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return false
	}
	if res.StatusCode/100 != 2 {
		return false
	}
	return true
}
