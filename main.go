package main

import (
	"os"

	"github.com/andersfylling/disgord"
	"github.com/aplombomb/boombot/config"
	"github.com/aplombomb/boombot/discord"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.ErrorLevel,
}

func main() {
	creds, err := config.GetSecrets()
	if err != nil {
		log.Fatalf("Error retrieving secrets: %+v", err)
	}

	client := disgord.New(disgord.Config{
		ProjectName: "BoomBot",
		BotToken:    creds.BotToken,
		Logger:      log,
		RejectEvents: []string{
			disgord.EvtTypingStart,
			disgord.EvtGuildMemberAdd,
			disgord.EvtPresenceUpdate,
			disgord.EvtGuildMemberUpdate,
			disgord.EvtGuildMemberRemove,
		},
		Presence: &disgord.UpdateStatusPayload{
			Game: &disgord.Activity{
				Name: "Copyrighted Music",
				Type: 1,
			},
		},
		Cache: &disgord.CacheNop{},
	})

	discord.BotRun(client, "&", "187298550193979392", creds.YoutubeToken)
}
