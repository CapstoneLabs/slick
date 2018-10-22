package mooder

import (
	"math/rand"
	"time"

	"github.com/CapstoneLabs/slick"
)

type Mooder struct {
	bot *slick.Bot
}

func init() {
	slick.RegisterPlugin(&Mooder{})
}

func (mooder *Mooder) InitPlugin(bot *slick.Bot) {
	mooder.bot = bot
	go mooder.SetupMoodChanger()
}

func (mooder *Mooder) SetupMoodChanger() {
	bot := mooder.bot
	for {
		time.Sleep(10 * time.Second)
		newMood := slick.Happy

		rand.Seed(time.Now().UTC().UnixNano())

		happyChances := rand.Int() % 10
		if happyChances > 6 {
			newMood = slick.Hyper
		}

		bot.Mood = newMood

		//bot.SendToChannel(bot.Config.GeneralChannel, bot.WithMood("I'm quite happy today.", "I can haz!! It's going to be a great one today!!"))

		select {
		case <-slick.AfterNextWeekdayTime(time.Now(), time.Monday, 12, 0):
		case <-slick.AfterNextWeekdayTime(time.Now(), time.Tuesday, 12, 0):
		case <-slick.AfterNextWeekdayTime(time.Now(), time.Wednesday, 12, 0):
		case <-slick.AfterNextWeekdayTime(time.Now(), time.Thursday, 12, 0):
		case <-slick.AfterNextWeekdayTime(time.Now(), time.Friday, 12, 0):
		}
	}
}
