package discord

import (
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

// Delegate evaluates commands and sends them to be processed
func (cec *CommandEventClient) Delegate() {
	cmd, _ := cec.DisectCommand()
	switch cmd {
	case "help", "h", "?", "wtf":
		hcc := NewHelpCommandClient(cec.data, cec.disgordClient)
		hcc.SendHelpMsg()
	case "play":
		// cec.data.Author
	default:

		uc := NewUnknownCommandClient(cec.data, cec.disgordClient)

		// err := Unknown(data.Message, client)

		uc.RespondToChannel()
	}

}

// DisectCommand returns the used command and all extraneous arguments
func (cec *CommandEventClient) DisectCommand() (string, []string) {
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
