package main

import (
	"fmt"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/aplombomb/boombot/config"
	"github.com/aplombomb/boombot/discord"
	"github.com/sirupsen/logrus"
)

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.ErrorLevel,
}

func main() {

	// Get the config from config.json
	conf := config.Retrieve("./config/config.json")
	fmt.Printf("\naccess key id: %+v\n", os.Getenv("AWS_ACCESS_KEY_ID"))
	// Fetch auth tokens from SecretsManager
	creds, err := config.GetSecrets()
	if err != nil {
		log.Fatalf("Error retrieving secrets: %+v", err)
	}

	client := disgord.New(disgord.Config{
		ProjectName: "BoomBot",
		BotToken:    creds.BotToken,
		Logger:      log,
		RejectEvents: []string{
			// rarely used, and causes unnecessary spam
			disgord.EvtTypingStart,
			// these require special privilege
			// https://discord.com/developers/docs/topics/gateway#privileged-intents
			disgord.EvtPresenceUpdate,
			disgord.EvtGuildMemberAdd,
			disgord.EvtGuildMemberUpdate,
			disgord.EvtGuildMemberRemove,
		},
		Presence: &disgord.UpdateStatusPayload{
			Game: &disgord.Activity{
				Name: "mewzek",
				// Type: 4,
				// Emoji: &disgord.ActivityEmoji{
				// 	Name:     "pepeSadJam",
				// 	ID:       779044106294263848,
				// 	Animated: true,
				// },
			},
		},
		Cache: &disgord.CacheNop{},
	})

	discord.BotRun(client, conf, creds)
}
