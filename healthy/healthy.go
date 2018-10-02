package healthy

import (
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlab.com/capstonemetering/cerberus/slick"
)

// Hipbot Plugin
type Healthy struct {
	urls []string
}

func init() {
	slick.RegisterPlugin(&Healthy{})
}

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

// Handler
func (healthy *Healthy) ChatHandler(listen *slick.Listener, msg *slick.Message) {
	log.Println("Health check. Requested by", msg.FromUser.Name)
	msg.Reply(healthy.CheckAll())
}

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
