package discord

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

//RespondToVoiceChannelUpdate contains logic for handling the voiceChannelUpdate event
func RespondToVoiceChannelUpdate(s disgord.Session, data *disgord.VoiceStateUpdate) {
	fmt.Printf("User %+v just joined the %+v voice chat", data.UserID, data.ChannelID)
}
