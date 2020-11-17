package discord

import (
	"context"
	"flag"
	"fmt"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/aplombomb/boombot/config"
)

// CmdArguments represents the arguments entered by the user after a command
type CmdArguments []string
type msgEvent disgord.Message

// TO-DO Get rid of these global variables
var ctx = context.Background()

var disgordGlobalClient *disgord.Client
var ytService *youtube.Service
var session disgord.Session
var conf config.ConfJSONStruct
var query = flag.String("query", "Google", "Search term")

// init the client
// var client = disgord.New(disgord.Config{BotToken: os.Getenv("BOOMBOT_TOKEN")})

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
	// sets the config for the whole disc package
	conf = cf

	disgordGlobalClient = client

	ytService, _ = youtube.NewService(ctx, option.WithAPIKey(creds.YoutubeToken))

	// filter incomming messages & set the prefix
	filter, _ := std.NewMsgFilter(ctx, client)
	filter.SetPrefix(cf.Prefix)
	//create a handler and bind it to new command events
	// go client.On(disgord.EvtMessageCreate,
	// 	filter.NotByBot,
	// 	filter.HasPrefix,
	// 	std.CopyMsgEvt,
	// 	filter.StripPrefix,

	// 	RespondToCommand,
	// )

	//Bind a handler to new message reactions
	// go client.On(disgord.EvtMessageReactionAdd, RespondToReaction)

	//Bind a handler to voice channel update events
	// go client.On(disgord.EvtVoiceStateUpdate, RespondToVoiceChannelUpdate)

	//Bind a handler to message events
	// go client.On(disgord.EvtMessageCreate, RespondToMessage)

	// The Gateway handler will replace the on handler once disgord becomse more stable
	// Keeping this here until that day comes
	client.Gateway().WithMiddleware(filter.NotByBot, filter.HasPrefix, std.CopyMsgEvt, filter.StripPrefix).MessageCreate(RespondToCommand)
	client.Gateway().MessageReactionAdd(RespondToReaction)
	client.Gateway().VoiceStateUpdate(RespondToVoiceChannelUpdate)
	client.Gateway().MessageCreate(RespondToMessage)

	fmt.Println("BoomBot is running")

	defer client.Gateway().StayConnectedUntilInterrupted()

}

func deleteMessage(resp *disgord.Message, sleep time.Duration, client *disgord.Client) {
	time.Sleep(sleep)

	fmt.Printf("\nDeleting message '%+v' by user %+v \n", resp.Content, resp.Author.Username)

	channel := client.Channel(resp.ChannelID)
	msgQueryBuilder := channel.Message(resp.ID)
	// if err != nil {
	// 	fmt.Printf("\nERROR DELETING MESSAGE: %+v\n", err)
	// }

	msgQueryBuilder.Delete()

}
