package discord

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

// Using this for access to the global clients FOR NOW as passing it through the handlers has proven tricky
// TO-DO find a solution to get rid of the global variables, including the client

// RespondToCommand delegates actions when commands are issued
func RespondToCommand(s disgord.Session, data *disgord.MessageCreate) {
	cec := NewCommandEventClient(data.Message, disgordGlobalClient)
	command, _ := cec.DisectCommand()

	fmt.Printf("\nvUserID: %+v\n", data.Message.Author.ID)

	user := data.Message.Author

	fmt.Printf("Command %+v by user %+v | %+v\n", command, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	switch command {
	case "play":
		go deleteMessage(data.Message, 20*time.Second, disgordGlobalClient)
	default:
		cec.Delegate()
	}

}

// RespondToMessage delegates actions when messages are created
func RespondToMessage(s disgord.Session, data *disgord.MessageCreate) {
	user := data.Message.Author

	fmt.Printf("Message %+v by user %+v | %+v\n", data.Message.Content, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	mec := NewMessageEventClient(data.Message, disgordGlobalClient)
	err := mec.FilterNonModLinks()
	if err != nil {
		fmt.Printf("\nError filtering non-mod link: %+v\n", err)
	}
}

// RespondToReaction delegates actions when reactions are added to messages
func RespondToReaction(s disgord.Session, data *disgord.MessageReactionAdd) {
	userQueryBuilder := disgordGlobalClient.User(data.UserID)

	user, err := userQueryBuilder.Get()

	if err != nil {
		fmt.Printf("\nError getting user: %+v\n", err)
	}
	// fmt.Printf("Message reaction %+v by user %+v | %+v\n", data.PartialEmoji.Name, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	rec := NewReactionEventClient(data.PartialEmoji, data.UserID, data.ChannelID, data.MessageID, disgordGlobalAPI)
	msg, err := rec.GenerateModResponse()
	if err != nil {
		fmt.Printf("\nError generating mod reaction response: %+v\n", err)
	}
	//TO-DO Sending the dm here as opposed having it sent via GenerateModResponse for testing purposes
	// Using it here at least allows me to get full coverage of the reactions logic
	// The SendMsg method of disgord.User requires session arg which has proven difficult to mock
	if msg != nil {
		user.SendMsg(ctx, s, msg)
	}
}

// RespondToVoiceChannelUpdate updates the server's voice channel cache every time an update is emitted
func RespondToVoiceChannelUpdate(s disgord.Session, data *disgord.VoiceStateUpdate) {
	globalQueue.UpdateVoiceCache(data.ChannelID, data.UserID)
}
