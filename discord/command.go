package discord

import (
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"

	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
	youtubeiface "github.com/aplombomb/boombot/youtube/ifaces"
)

// CommandEventClient contains the data for all command processing
type CommandEventClient struct {
	data          *disgord.Message
	disgordClient disgordiface.DisgordClientAPI
	youtubeClient youtubeiface.YoutubeSearchServiceAPI
}

// NewCommandEventClient returns a pointer to a new CommandEventClient
func NewCommandEventClient(data *disgord.Message, disgordClient disgordiface.DisgordClientAPI, youtubeClient youtubeiface.YoutubeSearchServiceAPI) *CommandEventClient {
	return &CommandEventClient{
		data:          data,
		disgordClient: disgordClient,
		youtubeClient: youtubeClient,
	}
}

//RespondToCommand handles all messages that begin with the configured prefix
func (cec *CommandEventClient) RespondToCommand() {
	cmd, _ := cec.DisectCommand(cec.data)

	switch cmd {
	case "help", "h", "?", "wtf":
		fmt.Println(cec.data.Content)
		hcc := NewHelpCommandClient(cec.data, cec.disgordClient)

		hcc.SendHelpMsg()
	case "play":

		// init the Youtube client here for test coverage's sake | will find another home for this later
		// ctx := context.Background()

		// if inVoice := ytClient.VerifyVoiceChat(s); inVoice == false {
		// 	fmt.Println("User is not in voice channel")
		// } else {
		// 	fmt.Println("User is in voice channel")
		// }

		// yt.PrintIt(ytClient)
	default:

		uc := NewUnknownCommandClient(cec.data, cec.disgordClient)

		// err := Unknown(data.Message, client)

		uc.RespondToChannel()
	}

}

// DisectCommand returns the used command and all extraneous arguments
func (cec *CommandEventClient) DisectCommand(data *disgord.Message) (string, []string) {
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
