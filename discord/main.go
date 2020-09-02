package discord

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/andersfylling/snowflake/v4"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/aplombomb/boombot/config"
)

// CmdArguments represents the arguments entered by the user after a command
type CmdArguments []string
type msgEvent disgord.Message

type AdminReaction struct {
	userID    snowflake.Snowflake
	channelID snowflake.Snowflake
	emoji     string
}

// Global Variables to ease working with client/sesion etc
var ctx = context.Background()
var client *disgord.Client
var session disgord.Session
var conf config.ConfJSONStruct

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

	//create a handler and bind it to new message events
	go client.On(disgord.EvtMessageCreate,
		filter.NotByBot,
		filter.HasPrefix,
		std.CopyMsgEvt,
		filter.StripPrefix,

		respondToMessage,
	)

	//Bind a handler to new message reactions
	go client.On(disgord.EvtMessageReactionAdd,

		respondToReaction,
	)

	go client.On(disgord.EvtVoiceStateUpdate,

		respondToVoiceChannelJoin,
	)

	fmt.Println("The bot is currently running")
}

func respondToMessage(s disgord.Session, data *disgord.MessageCreate) {
	cmd, args := ParseMessage(data)

	switch cmd {
	case "help", "h", "?", "wtf":
		help(data, args)
	}

}

func respondToReaction(s disgord.Session, data *disgord.MessageReactionAdd) {
	fmt.Printf("Name: %+v\nChannelID: %+v\nUserID: %+v\n", data.PartialEmoji.Name, data.ChannelID, data.UserID)
	reaction := ParseReaction(data)
	seenReaction := &AdminReaction{
		userID:    321044596476084235,
		channelID: 734986357583380510,
		emoji:     "👀",
	}
	fmt.Printf("Reaction: %+v\n", reaction)
	//if the reaction has been added to a message by me with the eye emoji
	//in the mod requests channel, send a message to the member that
	//suggested the mod that i have seen their suggestion
	if reflect.DeepEqual(reaction, seenReaction) {
		seenMsg := disgord.Message{
			Content: "Bomb has seen your mod recommendation",
		}
		message, _ := client.GetMessage(ctx, data.ChannelID, data.MessageID)
		message.Author.SendMsg(ctx, s, &seenMsg)

	}
}

func respondToVoiceChannelJoin(s disgord.Session, data *disgord.VoiceStateUpdate) {
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
func ParseReaction(data *disgord.MessageReactionAdd) *AdminReaction {
	return &AdminReaction{
		userID:    data.UserID,
		channelID: data.ChannelID,
		emoji:     data.PartialEmoji.Name,
	}
}

func deleteMessage(resp *disgord.Message, sleep time.Duration) {
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
