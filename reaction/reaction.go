package reaction

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
)

//Reaction defines the structure of needed reaction data
type Reaction struct {
	userID    snowflake.Snowflake
	channelID snowflake.Snowflake
	emoji     string
	action    string
}

//ModReactions defines the data structure for reactions that
//the bot should respond to inside the 'Mod Requests' channel
type ModReactions struct {
	SeenReactions     []Reaction
	AcceptedReactions []Reaction
	RejectedReactions []Reaction
}

//New returns a new reaction struct
func New(uID snowflake.Snowflake, chID snowflake.Snowflake, emoji string) *Reaction {
	return &Reaction{
		userID:    uID,
		channelID: chID,
		emoji:     emoji,
	}
}

//RespondToReaction processes\Delegates reaction events
func RespondToReaction(ctx context.Context, client *disgord.Client, s disgord.Session, data *disgord.MessageReactionAdd, reactionPool *ModReactions) {

	fmt.Println(reactionPool)
	reactionEvent := New(data.UserID, data.ChannelID, data.PartialEmoji.Name)

	//Loop through valid seen reactions and check for a match
	//TODO-These loops need to be consolidated into a single function
	for _, currentSeenReaction := range reactionPool.SeenReactions {
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
			go DeleteMessage(ctx, client, message, 3600)
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
			go DeleteMessage(ctx, client, message, 3600)
			break
		}
	}

}

//HydrateModReactions creates the ModReactions struct used
//for the bot to respond to specific reactions
func (mr *ModReactions) HydrateModReactions(seenEmojis []string, acceptedEmojis []string, rejectedEmojis []string) *ModReactions {
	sr := []Reaction{}
	ar := []Reaction{}
	rr := []Reaction{}

	for _, emoji := range seenEmojis {
		fmt.Println(emoji)
		sr = append(sr, Reaction{
			userID:    321044596476084235,
			channelID: 734986357583380510,
			emoji:     emoji,
		})
	}
	for _, emoji := range acceptedEmojis {
		sr = append(ar, Reaction{
			userID:    321044596476084235,
			channelID: 734986357583380510,
			emoji:     emoji,
		})
	}
	for _, emoji := range rejectedEmojis {
		sr = append(rr, Reaction{
			userID:    321044596476084235,
			channelID: 734986357583380510,
			emoji:     emoji,
		})
	}

	return &ModReactions{
		SeenReactions:     sr,
		AcceptedReactions: ar,
		RejectedReactions: rr,
	}

}

//DeleteMessage deletes the message after the specified time
func DeleteMessage(ctx context.Context, client *disgord.Client, resp *disgord.Message, sleep time.Duration) {
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
