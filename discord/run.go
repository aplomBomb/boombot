package discord

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/aplombomb/boombot/config"
)

// CmdArguments represents the arguments entered by the user after a command
type CmdArguments []string
type msgEvent disgord.Message

// Global Variables to ease working with client/sesion etc
var ctx = context.Background()

// var client *disgord.Client
var session disgord.Session
var conf config.ConfJSONStruct

// init the client
var client = disgord.New(disgord.Config{BotToken: os.Getenv("BOOMBOT_TOKEN")})

//Version of BoomBot
const Version = "v0.0.0-alpha"

func init() {
	//BoomBot cli logo
	fmt.Printf(`
	▄▄▄▄·             • ▌ ▄ ·.  ▄▄▄▄      ▄▄▄▄▄▄▄
	▐█ ▀█▪ ▄█▀▄  ▄█▀▄ ·██ ▐███▪▐█ ▀█▪ ▄█▀▄ •██
	▐█▀▀█▄▐█▌.▐▌▐█▌.▐▌▐█ ▌▐▌▐█·▐█▀▀█▄▐█▌.▐▌ ▐█.▪
	██▄▪▐█▐█▌.▐▌▐█▌.▐▌██ ██▌▐█▌██▄▪▐█▐█▌.▐▌ ▐█▌·
	·▀▀▀▀  ▀█▄▀▪ ▀█▄▀▪▀▀  █▪▀▀▀·▀▀▀▀  ▀█▄▀▪ ▀▀▀ %-16s\/`+"\n\n", Version)
}

// BotRun | Start the bot and handle events
func BotRun(cf config.ConfJSONStruct) {
	// sets the config for the whole disc package
	conf = cf

	// stay connected to discord
	defer client.StayConnectedUntilInterrupted(ctx)

	// filter incomming messages & set the prefix
	filter, _ := std.NewMsgFilter(ctx, client)
	filter.SetPrefix(cf.Prefix)

	//create a handler and bind it to new command events
	go client.On(disgord.EvtMessageCreate,
		filter.NotByBot,
		filter.HasPrefix,
		std.CopyMsgEvt,
		filter.StripPrefix,

		RespondToCommand,
	)

	//Bind a handler to new message reactions
	go client.On(disgord.EvtMessageReactionAdd, RespondToReaction)

	//Bind a handler to voice channel update events
	go client.On(disgord.EvtVoiceStateUpdate, RespondToVoiceChannelUpdate)

	//Bind a handler to message events
	go client.On(disgord.EvtMessageCreate, RespondToMessage)

	fmt.Println("BoomBot is running")
}

func deleteMessage(resp *disgord.Message, sleep time.Duration, client *disgord.Client) {
	time.Sleep(sleep)

	err := client.DeleteMessage(
		ctx,
		resp.ChannelID,
		resp.ID,
	)
	if err != nil {
		fmt.Println("error deleting message :", err)
	}
}
