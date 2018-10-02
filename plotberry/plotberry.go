package plotberry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"math"
	"net/http"
	"time"

	"gitlab.com/capstonemetering/cerberus/slick"
)

type PlotBerry struct {
	bot        *slick.Bot
	totalUsers int
	pingTime   time.Duration
	celebrated bool
}

type TotalUsers struct {
	Plotberries int `json:"plotberries"`
}

func init() {
	slick.RegisterPlugin(&PlotBerry{})
}

func (plotberry *PlotBerry) InitPlugin(bot *slick.Bot) {

	plotberry.bot = bot
	plotberry.celebrated = true
	plotberry.pingTime = 10 * time.Second
	plotberry.totalUsers = 100001

	statchan := make(chan TotalUsers, 100)

	go plotberry.launchWatcher(statchan)
	go plotberry.launchCounter(statchan)

	bot.Listen(&slick.Listener{
		MessageHandlerFunc: plotberry.ChatHandler,
	})
}

func (plotberry *PlotBerry) ChatHandler(listen *slick.Listener, msg *slick.Message) {
	if msg.MentionsMe && msg.Contains("how many user") {
		msg.Reply(fmt.Sprintf("We got %d users!", plotberry.totalUsers))
	}
	return
}

func getplotberry() (*TotalUsers, error) {

	var data TotalUsers

	resp, err := http.Get("https://plot.ly/v0/plotberries")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (plotberry *PlotBerry) launchWatcher(statchan chan TotalUsers) {

	for {
		time.Sleep(plotberry.pingTime)

		data, err := getplotberry()

		if err != nil {
			log.Print(err)
			continue
		}

		if data.Plotberries != plotberry.totalUsers {
			statchan <- *data
		}

		plotberry.totalUsers = data.Plotberries
	}
}

func (plotberry *PlotBerry) launchCounter(statchan chan TotalUsers) {

	finalcountdown := 100000

	for data := range statchan {

		totalUsers := data.Plotberries

		mod := math.Mod(float64(totalUsers), 50) == 0
		rem := finalcountdown - totalUsers

		if plotberry.celebrated {
			continue
		}

		if mod || (rem <= 10) {
			var msg string

			if rem == 10 {
				msg = fmt.Sprintf("@all %d users till the finalcountdown!", rem)
			} else if rem == 9 {
				msg = fmt.Sprintf("%d users!", rem)
			} else if rem == 8 {
				msg = fmt.Sprintf("and %d", rem)
			} else if rem == 7 {
				msg = fmt.Sprintf("we're at %d users. %d users till Mimosa time!\n", totalUsers, rem)
			} else if rem == 6 {
				msg = fmt.Sprintf("%d...", rem)
			} else if rem == 5 {
				msg = fmt.Sprintf("@all %d users\n I'm a freaky proud robot!", rem)
			} else if rem == 4 {
				msg = fmt.Sprintf("%d users till finalcountdown!", rem)
			} else if rem == 3 {
				msg = fmt.Sprintf("%d... \n", rem)
			} else if rem == 2 {
				msg = fmt.Sprintf("%d more! humpa humpa\n", rem)
			} else if rem == 1 {
				plotberry.bot.SendToChannel(plotberry.bot.Config.GeneralChannel, fmt.Sprintf("%d users until 100000.\nYOU'RE ALL MAGIC!", rem))
				msg = "https://31.media.tumblr.com/3b74abfa367a3ed9a2cd753cd9018baa/tumblr_miul04oqog1qkp8xio1_400.gif"
			} else if rem <= 0 {
				msg = fmt.Sprintf("@all FINALCOUNTDOWN!!!\n We're at %d user signups!!!!! My human compatriots, taking an idea to a product with 100,000 users is an achievement few will experience in their life times. Reflect, humans, on your hard work and celebrate this success. You deserve it, and remember, Plot On!", totalUsers)
				plotberry.celebrated = true
			} else {
				msg = fmt.Sprintf("We are at %d total user signups!", totalUsers)
			}

			plotberry.bot.SendToChannel(plotberry.bot.Config.GeneralChannel, msg)
		}
	}

}
