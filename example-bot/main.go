package main

import (
	"flag"
	"os"

	"gitlab.com/capstonemetering/cerberus/slick"
	_ "gitlab.com/capstonemetering/cerberus/slick/bugger"
	_ "gitlab.com/capstonemetering/cerberus/slick/deployer"
	_ "gitlab.com/capstonemetering/cerberus/slick/faceoff"
	_ "gitlab.com/capstonemetering/cerberus/slick/funny"
	_ "gitlab.com/capstonemetering/cerberus/slick/healthy"
	_ "gitlab.com/capstonemetering/cerberus/slick/hooker"
	_ "gitlab.com/capstonemetering/cerberus/slick/mooder"
	_ "gitlab.com/capstonemetering/cerberus/slick/plotberry"
	_ "gitlab.com/capstonemetering/cerberus/slick/recognition"
	_ "gitlab.com/capstonemetering/cerberus/slick/standup"
	_ "gitlab.com/capstonemetering/cerberus/slick/todo"
	_ "gitlab.com/capstonemetering/cerberus/slick/totw"
	_ "gitlab.com/capstonemetering/cerberus/slick/web"
	_ "gitlab.com/capstonemetering/cerberus/slick/webauth"
	_ "gitlab.com/capstonemetering/cerberus/slick/webutils"
	_ "gitlab.com/capstonemetering/cerberus/slick/wicked"
)

var configFile = flag.String("config", os.Getenv("HOME")+"/.slick.conf", "config file")

func main() {
	flag.Parse()

	bot := slick.New(*configFile)

	bot.Run()
}
