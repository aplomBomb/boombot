package discord

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
)

var (
	seenEmojis = []*disgord.Emoji{
		{
			Name: "👀",
		},
		{
			Name: "eyes",
		},
		{
			Name: "monkaEyesZoom",
		},
		{
			Name: "eyesFlipped",
		},
		{
			Name: "freakouteyes",
		},
		{
			Name: "monkaUltraEyes",
		},
		{
			Name: "PepeHmm",
		},
	}
	acceptedEmojis = []*disgord.Emoji{
		{
			Name: "✅",
		},
		{
			Name: "check",
		},
		{
			Name: "👍",
		},
		{
			Name: "ablobyes",
		},
		{
			Name: "Check",
		},
		{
			Name: "seemsgood",
		},
	}
	rejectedEmojis = []*disgord.Emoji{
		{
			Name: "🚫",
		},
		{
			Name: "no",
		},
		{
			Name: "steve_nope",
		},
		{
			Name: "❌",
		},
		{
			Name: "xmark",
		},
	}
)

//AdminReaction defines the structure of needed reaction data
type AdminReaction struct {
	userID    snowflake.Snowflake
	channelID snowflake.Snowflake
	emoji     *disgord.Emoji
}

// ReactionEventClient defines contextual data regarding a message react event
type ReactionEventClient struct {
	emoji         *disgord.Emoji
	uID           disgord.Snowflake
	chID          disgord.Snowflake
	msgID         disgord.Snowflake
	disgordClient disgordiface.DisgordClientAPI
}

// NewReactionEventClient returns a pointer to a new ReactionEventClient
func NewReactionEventClient(emoji *disgord.Emoji, uID disgord.Snowflake, chID disgord.Snowflake, msgID disgord.Snowflake, disgordClient disgordiface.DisgordClientAPI) *ReactionEventClient {
	return &ReactionEventClient{
		emoji,
		uID,
		chID,
		msgID,
		disgordClient,
	}
}

// HandleJukeboxReact triggers the playback channels of the queue in response to user reaction
func (rec *ReactionEventClient) HandleJukeboxReact(queue *Queue) {
	if rec.uID != 860286976296878080 && rec.uID == queue.NowPlayingUID {
		switch rec.emoji.Name {
		case "\u26D4":
			go func() {
				queue.Stop <- true
				return
			}()
		case "\u267B":
			go func() {
				queue.Shuffle <- true
				return
			}()
		case "\u23F8":
			go func() {
				queue.Pause <- true
				return
			}()
		case "\u25B6":
			go func() {
				queue.Play <- true
				return
			}()
		case "\u23E9":
			go func() {
				queue.Next <- true
				return
			}()
		}
	}
}

//GenerateModResponse returns the applicable message response if reaction criteria are met
func (rec *ReactionEventClient) GenerateModResponse() (*disgord.Message, error) {
	// fmt.Printf("Name: %+v\nChannelID: %+v\nUserID: %+v\n", rec.emoji.Name, rec.chID, rec.uID)

	dm := disgord.Message{}

	reactionEvent := &AdminReaction{
		userID:    rec.uID,
		channelID: rec.chID,
		emoji:     rec.emoji,
	}

	seenReactions := createReactions(seenEmojis)
	acceptedReactions := createReactions(acceptedEmojis)
	rejectedReactions := createReactions(rejectedEmojis)

	//Loop through valid seen reactions and check for a match
	// TODO: These loops need to be consolidated into a single function
	for _, currentSeenReaction := range seenReactions {
		if reflect.DeepEqual(currentSeenReaction, reactionEvent) {
			url := ""
			modName := ""
			chQueryBuilder := rec.disgordClient.Channel(rec.chID)
			message, _ := chQueryBuilder.GetMessages(&disgord.GetMessagesParams{
				Around: rec.msgID,
			})
			msgFields := strings.Fields(message[0].Content)

			//snag the url and the mod name from the request
			for _, field := range msgFields {
				if strings.Contains(field, "https://www.curseforge.com/minecraft/mc-mods/") {
					url = field
					urlFields := strings.Split(url, "/")
					i := len(urlFields) - 1
					modName = urlFields[i]
				}
			}
			dm = disgord.Message{
				Embeds: []*disgord.Embed{
					{
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
			break
		}
	}
	//Loop through valid accepted reactions and check for a match
	for _, currentAcceptedReaction := range acceptedReactions {
		if reflect.DeepEqual(currentAcceptedReaction, reactionEvent) {
			url := ""
			modName := ""
			chQueryBuilder := rec.disgordClient.Channel(rec.chID)
			message, _ := chQueryBuilder.GetMessages(&disgord.GetMessagesParams{
				Around: rec.msgID,
			})
			msgFields := strings.Fields(message[0].Content)

			//snag the url and the mod name from the request
			for _, field := range msgFields {
				if strings.Contains(field, "https://www.curseforge.com/minecraft/mc-mods/") {
					url = field
					urlFields := strings.Split(url, "/")
					i := len(urlFields) - 1
					modName = urlFields[i]
				}
			}
			dm = disgord.Message{
				Embeds: []*disgord.Embed{
					{
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
			go deleteMessage(message[0], 3*time.Second, rec.disgordClient)
			break
		}
	}
	//Loop through valid rejected reactions and check for a match
	//Extract the mod name to include in the embedded dm to the user for context
	for _, currentRejectedReaction := range rejectedReactions {
		if reflect.DeepEqual(currentRejectedReaction, reactionEvent) {
			url := ""
			modName := ""
			chQueryBuilder := rec.disgordClient.Channel(rec.chID)
			message, _ := chQueryBuilder.GetMessages(&disgord.GetMessagesParams{
				Around: rec.msgID,
			})
			msgFields := strings.Fields(message[0].Content)

			//snag the url and the mod name from the request
			for _, field := range msgFields {
				if strings.Contains(field, "https://www.curseforge.com/minecraft/mc-mods/") {
					url = field
					urlFields := strings.Split(url, "/")
					i := len(urlFields) - 1
					modName = urlFields[i]
				}
			}

			dm = disgord.Message{
				Embeds: []*disgord.Embed{
					{
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
			go deleteMessage(message[0], 3*time.Second, rec.disgordClient)
			break
		}
	}
	return &dm, nil
}

//ParseReaction bundles up reaction data for easier comparison
func createReactions(emojis []*disgord.Emoji) []*AdminReaction {
	reactions := []*AdminReaction{}
	for _, emoji := range emojis {
		reactions = append(reactions, &AdminReaction{
			userID:    801961246429544469,
			channelID: 734986357583380510,
			emoji:     emoji,
		})
	}
	return reactions
}
