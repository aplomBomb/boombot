package discord

import (
	"github.com/andersfylling/disgord"
)

// Using this for access to the global clients FOR NOW as passing it through the handlers has proven tricky
// TO-DO find a solution to get rid of the global variables, including the client

// RespondToCommand delegates actions when commands are issued
func RespondToCommand(s disgord.Session, data *disgord.MessageCreate) {
	cec := NewCommandEventClient(data.Message, disgordGlobalClient, ytService.Search)
	cec.RespondToCommand()
}

// RespondToMessage delegates actions when messages are created
func RespondToMessage(s disgord.Session, data *disgord.MessageCreate) {
	mec := NewMessageEventClient(data.Message, disgordGlobalClient)
	mec.FilterNonModLinks()
}

// RespondToReaction delegates actions when reactions are added to messages
func RespondToReaction(s disgord.Session, data *disgord.MessageReactionAdd) {
	rec := NewReactionEventClient(data.PartialEmoji, data.UserID, data.ChannelID, data.MessageID, disgordGlobalClient)
	rec.RespondToReaction(s)
}

// RespondToVoiceChannelUpdate delegates actions when voice state events are triggered
func RespondToVoiceChannelUpdate(s disgord.Session, data *disgord.VoiceStateUpdate) {

}
