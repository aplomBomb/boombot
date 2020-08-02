package discord

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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
var client *disgord.Client
var session disgord.Session
var conf config.ConfJSONStruct

// BotRun the bot and handle events
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

	// create a handler and bind it to new message events
	go client.On(disgord.EvtMessageCreate,
		// middleware
		filter.NotByBot,    // ignore bot messages
		filter.HasPrefix,   // read original
		std.CopyMsgEvt,     // read & copy original
		filter.StripPrefix, // write copy

		// handler
		reply, // call reply func
	) // handles copy

	fmt.Println("The bot is currently running")
}

func reply(s disgord.Session, data *disgord.MessageCreate) {
	cmd, args := ParseMessage(data)

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

// ParseArgToSearch parses any arg to an int, if no int is entered, returns 0 as the result
func (args CmdArguments) ParseArgToSearch() query.CharSearchInput {
	id, err := strconv.Atoi(args[0])
	arg := strings.Join(args, " ")
	if err != nil && id != 0 {
		fmt.Println(err)
	}
	return query.CharSearchInput{ID: id, Name: arg}
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
