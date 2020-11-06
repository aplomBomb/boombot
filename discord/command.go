package discord

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/andersfylling/disgord"
	yt "github.com/aplombomb/boombot/Youtube"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
)

// CommandEventClient contains the data for all command processing
type CommandEventClient struct {
	data          *disgord.Message
	disgordClient disgordiface.DisgordClientAPI
}

// NewCommandEventClient returns a pointer to a new CommandEventClient
func NewCommandEventClient(data *disgord.Message, disgordClient disgordiface.DisgordClientAPI) *CommandEventClient {
	return &CommandEventClient{
		data:          data,
		disgordClient: disgordClient,
	}
}

//RespondToCommandTemp handles all messages that begin with the configured prefix
func (cec *CommandEventClient) RespondToCommandTemp(s disgord.Session, data *disgord.MessageCreate) {
	cmd, _ := cec.ParseCommand(data.Message)

	switch cmd {
	case "help", "h", "?", "wtf":
		fmt.Println(data.Message.Content)
		hcc := NewHelpCommandClient(data.Message, client)

		hcc.SendHelpMsg()
	case "play":

		// init the Youtube client here for test coverage's sake | will find another home for this later
		ctx := context.Background()

		youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(conf.YoutubeToken))

		if err != nil {
			fmt.Println(err)
		}
		ytClient, err := yt.New(youtubeService, data.Message.Content, data.Message.Author)

		if err != nil {
			log.Fatal("YT API ERROR: ", err)
		}

		fmt.Printf("\nYT Client Created: %+v\n\n\n", ytClient)

		// if inVoice := ytClient.VerifyVoiceChat(s); inVoice == false {
		// 	fmt.Println("User is not in voice channel")
		// } else {
		// 	fmt.Println("User is in voice channel")
		// }

		// yt.PrintIt(ytClient)
	default:

		uc := NewUnknownCommandClient(data.Message, client)

		// err := Unknown(data.Message, client)

		uc.RespondToChannel()
	}

}

// ParseCommand returns the used command and all extraneous arguments
func (cec *CommandEventClient) ParseCommand(data *disgord.Message) (string, []string) {
	var command string
	var args []string
	if len(cec.data.Content) > 0 {
		command = strings.ToLower(strings.Fields(cec.data.Content)[0])
		if len(cec.data.Content) > 1 {
			args = strings.Fields(cec.data.Content)[1:]
		}
	}
	return command, args
}
