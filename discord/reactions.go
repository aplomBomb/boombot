package discord

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
)

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
