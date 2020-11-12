package main

import (
	"log"

	"github.com/andersfylling/disgord"
	"github.com/aplombomb/boombot/config"
	"github.com/aplombomb/boombot/discord"
)

func main() {

	// Get the config from config.json
	conf := config.Retrieve("./config/config.json")

	creds, err := config.GetSecrets()
	if err != nil {
		log.Fatalf("Error retrieving secrets: %+v", err)
	}

	client := disgord.New(disgord.Config{
		BotToken: creds.BotToken,
		// Will use this in future disgord version once it actually works
		// Cache:    &disgord.CacheNop{},
	})

	discord.BotRun(client, conf, creds)
}
