package discord

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
)

// CmdArguments represents the arguments entered by the user after a command
type CmdArguments []string
type msgEvent disgord.Message

// TO-DO Get rid of these global variables
var ctx = context.Background()
var disgordGlobalAPI disgordiface.DisgordClientAPI
var disgordGlobalClient *disgord.Client
var globalGuild disgord.GuildQueryBuilder
var ytService *youtube.Service
var session disgord.Session
var globalQueue *Queue

// Version of BoomBot
const Version = "v1.0.0-alpha"

func init() {

	fmt.Printf(`
	▄▄▄▄·             • ▌ ▄ ·.  ▄▄▄▄      ▄▄▄▄▄▄▄
	▐█ ▀█▪ ▄█▀▄  ▄█▀▄ ·██ ▐███▪▐█ ▀█▪ ▄█▀▄ •██
	▐█▀▀█▄▐█▌.▐▌▐█▌.▐▌▐█ ▌▐▌▐█·▐█▀▀█▄▐█▌.▐▌ ▐█.▪
	██▄▪▐█▐█▌.▐▌▐█▌.▐▌██ ██▌▐█▌██▄▪▐█▐█▌.▐▌ ▐█▌·
	·▀▀▀▀  ▀█▄▀▪ ▀█▄▀▪▀▀  █▪▀▀▀·▀▀▀▀  ▀█▄▀▪ ▀▀▀ %-16s\/`+"\n\n", Version)
}

// BotRun | Start the bot and handle events
func BotRun(client *disgord.Client, prefix string, gID string, yk string) {
	queue := NewQueue(disgord.ParseSnowflakeString(gID))
	globalQueue = queue
	disgordGlobalClient = client
	gb := disgordGlobalClient.Guild(disgord.ParseSnowflakeString(gID))
	globalGuild = gb
	ytService, _ = youtube.NewService(ctx, option.WithAPIKey(yk))
	vlc := ytService.Videos.List([]string{"contentDetails", "snippet", "statistics"})
	filter, _ := std.NewMsgFilter(ctx, client)
	filter.SetPrefix(prefix)
	client.Gateway().WithMiddleware(filter.NotByBot, filter.HasPrefix, std.CopyMsgEvt, filter.StripPrefix).MessageCreate(RespondToCommand)
	client.Gateway().MessageReactionAdd(RespondToReaction)
	client.Gateway().VoiceStateUpdate(RespondToVoiceChannelUpdate)
	client.Gateway().MessageCreate(RespondToMessage)
	fmt.Println("BoomBot is running")
	go globalQueue.ListenAndProcessQueue(client, gb, vlc)
	go globalQueue.ManageJukebox(client)
	defer client.Gateway().StayConnectedUntilInterrupted()
}
