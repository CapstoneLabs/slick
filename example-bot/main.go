package main

import (
	"flag"
	"os"

	"github.com/CapstoneLabs/slick"
	_ "github.com/CapstoneLabs/slick/bugger"
	_ "github.com/CapstoneLabs/slick/deployer"
	_ "github.com/CapstoneLabs/slick/faceoff"
	_ "github.com/CapstoneLabs/slick/funny"
	_ "github.com/CapstoneLabs/slick/healthy"
	_ "github.com/CapstoneLabs/slick/hooker"
	_ "github.com/CapstoneLabs/slick/mooder"
	_ "github.com/CapstoneLabs/slick/plotberry"
	_ "github.com/CapstoneLabs/slick/recognition"
	_ "github.com/CapstoneLabs/slick/standup"
	_ "github.com/CapstoneLabs/slick/todo"
	_ "github.com/CapstoneLabs/slick/totw"
	_ "github.com/CapstoneLabs/slick/web"
	_ "github.com/CapstoneLabs/slick/webauth"
	_ "github.com/CapstoneLabs/slick/webutils"
	_ "github.com/CapstoneLabs/slick/wicked"
)

var configFile = flag.String("config", os.Getenv("HOME")+"/.slick.conf", "config file")

func main() {
	flag.Parse()

	bot := slick.New(*configFile)

	bot.Run()
}
