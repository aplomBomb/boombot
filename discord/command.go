package discord

import (
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"

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

// Delegate evaluates commands and sends them to be processed
func (cec *CommandEventClient) Delegate() {
	cmd, _ := cec.DisectCommand()
	switch cmd {
	case "help", "h", "?", "wtf":
		hcc := NewHelpCommandClient(cec.data, cec.disgordClient)
		hcc.SendHelpMsg()
	default:
		uc := NewUnknownCommandClient(cec.data, cec.disgordClient)
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

// Print the ID and title of each result in a list as well as a name that
// identifies the list. For example, print the word section name "Videos"
// above a list of video search results, followed by the video ID and title
// of each matching video.
func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}
