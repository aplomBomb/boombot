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
	user, err := disgordGlobalClient.GetUser(ctx, data.Message.Author.ID)
	if err != nil {
		fmt.Println("Failed to fetch user")
		user = &disgord.User{
			Username: "unknown",
		}
	}
	cec := NewCommandEventClient(data.Message, disgordGlobalClient, ytService.Search)
	command, _ := cec.DisectCommand()
	fmt.Printf("Command %+v by user %+v | %+v\n", command, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	cec.RespondToCommand()
}

// RespondToMessage delegates actions when messages are created
func RespondToMessage(s disgord.Session, data *disgord.MessageCreate) {
	user, err := disgordGlobalClient.GetUser(ctx, data.Message.Author.ID)
	if err != nil {
		fmt.Println("Failed to fetch user (probably a webhook message)")
		user = &disgord.User{
			Username: "unknown",
		}
	}
	fmt.Printf("Message %+v by user %+v | %+v\n", data.Message.Content, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	mec := NewMessageEventClient(data.Message, disgordGlobalClient)
	mec.FilterNonModLinks()
}

// RespondToReaction delegates actions when reactions are added to messages
func RespondToReaction(s disgord.Session, data *disgord.MessageReactionAdd) {
	// user, _ := disgordGlobalClient.GetUser(ctx, data.UserID)
	// fmt.Printf("Message reaction %+v by user %+v | %+v\n", data.PartialEmoji.Name, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	rec := NewReactionEventClient(data.PartialEmoji, data.UserID, data.ChannelID, data.MessageID, disgordGlobalClient, s)
	rec.RespondToReaction()
}

// RespondToVoiceChannelUpdate delegates actions when voice state events are triggered
func RespondToVoiceChannelUpdate(s disgord.Session, data *disgord.VoiceStateUpdate) {

}
