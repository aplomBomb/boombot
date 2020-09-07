package discord

import (
	"context"
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/aplombomb/boombot/config"
	"github.com/aplombomb/boombot/reaction"
)

// CmdArguments represents the arguments entered by the user after a command
type CmdArguments []string
type msgEvent disgord.Message

// Global Variables to ease working with client/sesion etc
var (
	ctx     = context.Background()
	client  *disgord.Client
	session disgord.Session
	conf    config.ConfJSONStruct
)

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

	// init the client
	client = disgord.New(disgord.Config{BotToken: cf.BotToken})

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
	go client.On(disgord.EvtMessageReactionAdd, ParseReaction)

	//Bind a handler to voice channel update events
	go client.On(disgord.EvtVoiceStateUpdate, RespondToVoiceChannelUpdate)

	//Bind a handler to message events
	go client.On(disgord.EvtMessageCreate, RespondToMessage)

	fmt.Println("BoomBot is running")
}

//RespondToMessage handles all messages created in the server
func RespondToMessage(s disgord.Session, data *disgord.MessageCreate) {
	//Per channel message event switch handler
	switch data.Message.ChannelID {
	case 734986357583380510:
		if strings.Contains(data.Message.Content, "https://www.curseforge.com/minecraft/mc-mods/") == false {
			message, _ := client.GetMessage(ctx, data.Message.ChannelID, data.Message.ID)
			go reaction.DeleteMessage(ctx, client, message, 1)
		}
	default:
		break
	}

}

//RespondToCommand handles all messages that begin with prefix
func RespondToCommand(s disgord.Session, data *disgord.MessageCreate) {
	cmd, args := ParseMessage(data)

	switch cmd {
	case "help", "h", "?", "wtf":
		help(data, args)
	default:
		unknown(data)
	}

}

//ParseReaction is being used to extract the data object from the EvtMessageReactionAdd event,
//not sure of a better way to do this
func ParseReaction(s disgord.Session, data *disgord.MessageReactionAdd) {
	modReactionPool := reaction.ModReactions{}
	modReactionPool.HydrateModReactions(conf.SeenEmojis, conf.AcceptedEmojis, conf.RejectedEmojis)
	reaction.RespondToReaction(ctx, client, s, data, &modReactionPool)

}

//RespondToVoiceChannelUpdate contains logic for handling the voiceChannelUpdate event
func RespondToVoiceChannelUpdate(s disgord.Session, data *disgord.VoiceStateUpdate) {
	fmt.Printf("User %+v just joined the %+v voice chat", data.UserID, data.ChannelID)
}

// ParseMessage parses the message into command / args
func ParseMessage(data *disgord.MessageCreate) (string, []string) {
	var command string
	var args []string
	if len(data.Message.Content) > 0 {
		command = strings.ToLower(strings.Fields(data.Message.Content)[0])
		if len(data.Message.Content) > 1 {
			args = strings.Fields(data.Message.Content)[1:]
		}
	}
	return command, args
}

//ParseReaction bundles up reaction data for easier comparison
// func createReactions(emojis []string, data *disgord.MessageReactionAdd) *AdminReactions {
// 	reactions := []*AdminReaction{}
// 	for _, emoji := range emojis {
// 		reactions = append(reactions, &AdminReaction{
// 			userID:    321044596476084235,
// 			channelID: 734986357583380510,
// 			emoji:     emoji,
// 		})
// 	}
// 	return &AdminReactions{
// 		Reactions: reactions,
// 	}
// }
