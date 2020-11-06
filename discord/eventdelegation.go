package discord

import (
	"github.com/andersfylling/disgord"
)

// Using this for access to the global client FOR NOW as passing it through the handlers has proven tricky
// TO-DO find a solution to get rid of the global variables, including the client

func RespondToCommand(s disgord.Session, data *disgord.MessageCreate) {
	cec := NewCommandEventClient(data.Message, client)

}

func RespondToMessage(s disgord.Session, data *disgord.MessageCreate) {
	mec := NewMessageEventClient(data.Message, client)

	mec.FilterNonModLinks()

}

func RespondToReaction(s disgord.Session, data *disgord.MessageReactionAdd) {

}

func RespondToVoiceChannelUpdate(s disgord.Session, data *disgord.VoiceStateUpdate) {

}
