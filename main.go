package main

import (
	"os"

	"github.com/andersfylling/disgord"
	"github.com/aplombomb/boombot/config"
	"github.com/aplombomb/boombot/discord"
)

func main() {

	// Get the config from config.json
	conf := config.Retrieve("./config/config.json")

	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("BOOMBOT_TOKEN"),
		// Cache:    &disgord.CacheNop{},
	})
	// run(client)

	discord.BotRun(client, conf)
}
