package discord

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/aplombomb/boombot/config"
	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
)

// CmdArguments represents the arguments entered by the user after a command
type CmdArguments []string
type msgEvent disgord.Message

// TO-DO Get rid of these global variables
var ctx = context.Background()
var disgordGlobalAPI disgordiface.DisgordClientAPI
var disgordGlobalClient *disgord.Client
var ytService *youtube.Service
var session disgord.Session
var conf config.ConfJSONStruct
var globalQueue *Queue

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
func BotRun(client *disgord.Client, cf config.ConfJSONStruct, creds *config.BoombotCreds) {
	conf = cf
	queue := NewQueue(disgord.ParseSnowflakeString(conf.GuildID))
	globalQueue = queue

	disgordGlobalClient = client
	ytService, _ = youtube.NewService(ctx, option.WithAPIKey(creds.YoutubeToken))
	filter, _ := std.NewMsgFilter(ctx, client)
	filter.SetPrefix(cf.Prefix)
	client.Gateway().WithMiddleware(filter.NotByBot, filter.HasPrefix, std.CopyMsgEvt, filter.StripPrefix).MessageCreate(RespondToCommand)
	client.Gateway().MessageReactionAdd(RespondToReaction)
	client.Gateway().VoiceStateUpdate(RespondToVoiceChannelUpdate)
	client.Gateway().MessageCreate(RespondToMessage)
	fmt.Println("BoomBot is running")
	go globalQueue.ListenAndProcessQueue(client)
	go globalQueue.ManageJukebox(client)
	defer client.Gateway().StayConnectedUntilInterrupted()
}
