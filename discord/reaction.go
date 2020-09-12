package discord

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
)

//Reaction defines the structure of needed reaction data
type reaction struct {
	userID    snowflake.Snowflake
	channelID snowflake.Snowflake
	emoji     string
	action    string
}

//reactionPool defines the data structure for reactions that
//the bot should respond to inside the 'Mod Requests' channel
type reactionPool struct {
	SeenReactions     []reaction
	AcceptedReactions []reaction
	RejectedReactions []reaction
}

//New returns a new reaction struct
func new(uID snowflake.Snowflake, chID snowflake.Snowflake, emoji string) *reaction {
	return &reaction{
		userID:    uID,
		channelID: chID,
		emoji:     emoji,
	}
}

//RespondToReaction processes\Delegates reaction events
func respondToReaction(ctx context.Context, client *disgord.Client, s disgord.Session, data *disgord.MessageReactionAdd, reactionPool *reactionPool) {
	fmt.Println("Responding...")
	fmt.Printf("ReactionPool: %+v", reactionPool)
	reactionEvent := new(data.UserID, data.ChannelID, data.PartialEmoji.Name)
	fmt.Printf("ReactionEvent: %+v", reactionEvent)

	//Loop through valid seen reactions and check for a match
	//TODO-These loops need to be consolidated into a single function
	for _, currentSeenReaction := range reactionPool.SeenReactions {
		if reflect.DeepEqual(currentSeenReaction, reactionEvent) {
			fmt.Println("It's a match!")
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
	for _, currentAcceptedReaction := range reactionPool.AcceptedReactions {
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
			go deleteMessage(ctx, client, message, 3600)
			break
		}
	}
	//Loop through valid rejected reactions and check for a match
	//Extract the mod name to include in the embedded dm to the user for context
	for _, currentRejectedReaction := range reactionPool.RejectedReactions {
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
			go deleteMessage(ctx, client, message, 3600)
			break
		}
	}

}

//HydratereactionPool creates the reactionPool struct used
//for the bot to respond to specific reactions
func (mr *reactionPool) hydrateReactionPool(seenEmojis []string, acceptedEmojis []string, rejectedEmojis []string) *reactionPool {
	sr := []reaction{}
	ar := []reaction{}
	rr := []reaction{}

	for _, emoji := range seenEmojis {
		// fmt.Println(emoji)
		sr = append(sr, reaction{
			userID:    321044596476084235,
			channelID: 734986357583380510,
			emoji:     emoji,
		})
	}
	for _, emoji := range acceptedEmojis {
		// fmt.Println(emoji)
		sr = append(ar, reaction{
			userID:    321044596476084235,
			channelID: 734986357583380510,
			emoji:     emoji,
		})
	}
	for _, emoji := range rejectedEmojis {
		// fmt.Println(emoji)
		sr = append(rr, reaction{
			userID:    321044596476084235,
			channelID: 734986357583380510,
			emoji:     emoji,
		})
	}

	return &reactionPool{
		SeenReactions:     sr,
		AcceptedReactions: ar,
		RejectedReactions: rr,
	}

}
