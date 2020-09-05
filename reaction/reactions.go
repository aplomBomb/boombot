package reaction

import (
	"log"

	"github.com/andersfylling/snowflake/v4"
)

//Reaction defines the structure of needed reaction data
type Reaction struct {
	userID    snowflake.Snowflake
	channelID snowflake.Snowflake
	emoji     string
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

//HydrateModReactions creates the ModReactions struct used
//for the bot to respond to specific reactions
func HydrateModReactions(emoji string, reactionType string) {
	types := [3]string{"seen", "accepted", "rejected"}
	match := false

	for _, currentType := range types {
		if reactionType == currentType {
			match = true
		}
	}
	if match == false {
		log.Fatal("You must specify the kind of reaction you are trying to create!\n Options: 'seen', 'accepted', 'rejected'")
	}

}
