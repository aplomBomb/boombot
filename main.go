package main

import (
	"github.com/aplombomb/boombot/config"
	"github.com/aplombomb/boombot/discord"
)

func main() {

	// Get the config from config.json
	conf := config.Retrieve("./config/config.json")

	discord.BotRun(conf)
}
