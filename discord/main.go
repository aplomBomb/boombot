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

//AdminReaction defines the structure of needed reaction data
type AdminReaction struct {
	userID    snowflake.Snowflake
	channelID snowflake.Snowflake
	emoji     string
}

//AdminReactions contains slice of AdminReaction
type AdminReactions struct {
	Reactions []*AdminReaction
}

// Global Variables to ease working with client/sesion etc
var ctx = context.Background()
var client *disgord.Client
var session disgord.Session
var conf config.ConfJSONStruct

var (
	seenEmojis = []string{
		"👀",
		"eyes",
		"monkaEyesZoom",
		"eyesFlipped",
		"freakouteyes",
		"monkaUltraEyes",
		"PepeHmm",
	}
	acceptedEmojis = []string{
		"✅",
		"check",
		"👍",
		"ablobyes",
		"Check",
		"seemsgood",
	}
	rejectedEmojis = []string{
		"🚫",
		"no",
		"steve_nope",
		"❌",
		"xmark",
		"🇽",
	}
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

	//create a handler and bind it to new message events
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

	fmt.Println("BoomBot is running")
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

//RespondToReaction contains logic for handling the reaction add event
func RespondToReaction(s disgord.Session, data *disgord.MessageReactionAdd) {
	fmt.Printf("Name: %+v\nChannelID: %+v\nUserID: %+v\n", data.PartialEmoji.Name, data.ChannelID, data.UserID)

	reactionEvent := &AdminReaction{
		userID:    data.UserID,
		channelID: data.ChannelID,
		emoji:     data.PartialEmoji.Name,
	}

	seenReactions := createReactions(seenEmojis, data)
	acceptedReactions := createReactions(acceptedEmojis, data)
	rejectedReactions := createReactions(rejectedEmojis, data)

	//Loop through valid seen reactions and check for a match
	//TODO-These loops need to be consolidated into a single function
	for _, currentSeenReaction := range seenReactions.Reactions {
		if reflect.DeepEqual(currentSeenReaction, reactionEvent) {
			url := ""
			modName := ""
			message, _ := client.GetMessage(ctx, data.ChannelID, data.MessageID)
			msgFields := strings.Fields(message.Content)

			//snag the url and the mod name from the request
			for _, field := range msgFields {
				if strings.Contains(field, "https://www.curseforge.com/minecraft/mc-mods/") {
					url = field
					urlFields := strings.Split(url, "/")
					i := len(urlFields) - 1
					modName = urlFields[i]
				}
			}
			dm := disgord.Message{
				Embeds: []*disgord.Embed{
					&disgord.Embed{
						Title:       fmt.Sprintf("**Your request to add %s is being reviewed**", modName),
						URL:         url,
						Description: fmt.Sprintf("*Bomb is reviewing your request to add %s*", modName),
						Color:       0xcc0000,
						Footer: &disgord.EmbedFooter{
							Text:    "Sit tight partner!",
							IconURL: "https://cdn.discordapp.com/emojis/745396324215685201.gif?v=1",
						},
					},
				},
			}
			message.Author.SendMsg(ctx, s, &dm)
			break
		}
	}
	//Loop through valid accepted reactions and check for a match
	for _, currentAcceptedReaction := range acceptedReactions.Reactions {
		if reflect.DeepEqual(currentAcceptedReaction, reactionEvent) {
			url := ""
			modName := ""
			message, _ := client.GetMessage(ctx, data.ChannelID, data.MessageID)
			msgFields := strings.Fields(message.Content)

			//snag the url and the mod name from the request
			for _, field := range msgFields {
				if strings.Contains(field, "https://www.curseforge.com/minecraft/mc-mods/") {
					url = field
					urlFields := strings.Split(url, "/")
					i := len(urlFields) - 1
					modName = urlFields[i]
				}
			}
			dm := disgord.Message{
				Embeds: []*disgord.Embed{
					&disgord.Embed{
						Title:       fmt.Sprintf("**%s ACCEPTED!!**", modName),
						URL:         url,
						Description: fmt.Sprintf("*Bomb has added %s to the modpack! If the server breaks now, it's all your fault!*", modName),
						Color:       0xcc0000,
						Footer: &disgord.EmbedFooter{
							Text:    "Pervert Steve is always watching...",
							IconURL: "https://cdn.discordapp.com/emojis/681217726412488767.png?v=1",
						},
					},
				},
			}
			message.Author.SendMsg(ctx, s, &dm)
			go deleteMessage(message, 3600)
			break
		}
	}
	//Loop through valid rejected reactions and check for a match
	//Extract the mod name to include in the embedded dm to the user for context
	for _, currentRejectedReaction := range rejectedReactions.Reactions {
		if reflect.DeepEqual(currentRejectedReaction, reactionEvent) {
			url := ""
			modName := ""
			message, _ := client.GetMessage(ctx, data.ChannelID, data.MessageID)
			msgFields := strings.Fields(message.Content)

			//snag the url and the mod name from the request
			for _, field := range msgFields {
				if strings.Contains(field, "https://www.curseforge.com/minecraft/mc-mods/") {
					url = field
					urlFields := strings.Split(url, "/")
					i := len(urlFields) - 1
					modName = urlFields[i]
				}
			}

			dm := disgord.Message{
				Embeds: []*disgord.Embed{
					&disgord.Embed{
						Title:       fmt.Sprintf("**%s Rejected**", modName),
						URL:         url,
						Description: fmt.Sprintf("*Bomb has rejected your request to add %s*", modName),
						Color:       0xcc0000,
						Footer: &disgord.EmbedFooter{
							Text:    "You have brought much shame upon your famiry",
							IconURL: "https://cdn.discordapp.com/emojis/662170922580574258.gif?v=1",
						},
					},
				},
			}
			message.Author.SendMsg(ctx, s, &dm)
			go deleteMessage(message, 3600)
			break
		}
	}
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
func createReactions(emojis []string, data *disgord.MessageReactionAdd) *AdminReactions {
	reactions := []*AdminReaction{}
	for _, emoji := range emojis {
		reactions = append(reactions, &AdminReaction{
			userID:    321044596476084235,
			channelID: 734986357583380510,
			emoji:     emoji,
		})
	}
	return &AdminReactions{
		Reactions: reactions,
	}
}

func deleteMessage(resp *disgord.Message, sleep time.Duration) {
	time.Sleep(sleep * time.Second)

	err := client.DeleteMessage(
		ctx,
		resp.ChannelID,
		resp.ID,
	)
	if err != nil {
		fmt.Println("error deleting message :", err)
	}
}
