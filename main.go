package main

import (
	"github.com/andersfylling/disgord"
	"github.com/aplombomb/boombot/config"
	"github.com/aplombomb/boombot/discord"
)

func main() {

	// Get the config from config.json
	conf := config.Retrieve("./config/config.json")

	client := disgord.New(disgord.Config{
		BotToken: conf.BotToken,
		// Will use this in future disgord version once it actually works
		// Cache:    &disgord.CacheNop{},
	})

	discord.BotRun(client, conf)
}
