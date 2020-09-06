package reaction

import (
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
func RespondToReaction(s disgord.Session, data *disgord.MessageReactionAdd, pool ModReactions) {

}

//HydrateModReactions creates the ModReactions struct used
//for the bot to respond to specific reactions
func (mr ModReactions) HydrateModReactions(seenEmojis []string, acceptedEmojis []string, rejectedEmojis []string) *ModReactions {
	sr := []Reaction{}
	ar := []Reaction{}
	rr := []Reaction{}

	for _, emoji := range seenEmojis {
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
