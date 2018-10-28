package main

import (
	"flag"

	"github.com/CapstoneLabs/slick"
	_ "github.com/CapstoneLabs/slick/bugger"
	_ "github.com/CapstoneLabs/slick/faceoff"
	_ "github.com/CapstoneLabs/slick/funny"
	_ "github.com/CapstoneLabs/slick/healthy"
	_ "github.com/CapstoneLabs/slick/hooker"
	_ "github.com/CapstoneLabs/slick/mooder"
	_ "github.com/CapstoneLabs/slick/plotberry"
	_ "github.com/CapstoneLabs/slick/recognition"
	_ "github.com/CapstoneLabs/slick/standup"
	_ "github.com/CapstoneLabs/slick/todo"
	_ "github.com/CapstoneLabs/slick/web"
	_ "github.com/CapstoneLabs/slick/webauth"
	_ "github.com/CapstoneLabs/slick/webutils"
	_ "github.com/CapstoneLabs/slick/wicked"
)

// Specify an alternative config file. Slick searches the working
// directory and your home folder by default for a file called
// `config.json`, `config.yaml`, or `config.toml` if no config
// file is specified
var configFile = flag.String("config", "", "config file")

func main() {
	flag.Parse()

	bot := slick.New(*configFile)

	bot.Run()
}
