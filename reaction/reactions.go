package reaction

import (
	"github.com/andersfylling/snowflake/v4"
)

//Reaction defines the structure of needed reaction data
type Reaction struct {
	userID    snowflake.Snowflake
	channelID snowflake.Snowflake
	emoji     string
}

//Reactions contains slice of AdminReaction
type Reactions struct {
	Reactions []Reaction
}

//New returns a new reaction struct
func (react *Reaction) New(uID snowflake.Snowflake, chID snowflake.Snowflake, emoji string) *Reaction {
	return &Reaction{
		userID:    uID,
		channelID: chID,
		emoji:     emoji,
	}
}
