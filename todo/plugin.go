package todo

import (
	log "github.com/sirupsen/logrus"

	"github.com/boltdb/bolt"
	"gitlab.com/capstonemetering/cerberus/slick"
)

type Plugin struct {
	bot   *slick.Bot
	store Store
}

func init() {
	slick.RegisterPlugin(&Plugin{})
}

func (p *Plugin) InitPlugin(bot *slick.Bot) {
	p.bot = bot

	err := bot.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		log.Fatalln("Couldn't create the `todos` bucket")
	}

	p.store = &boltStore{db: bot.DB}
	p.listenTodo()
}
