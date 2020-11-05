package discord

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
	discordiface "github.com/aplombomb/boombot/discord/ifaces"
)

// HelpCommandClient contains the resources needed for handling help requests
type HelpCommandClient struct {
	data          *disgord.MessageCreate
	disgordClient discordiface.DisgordClientAPI
}

// func help(data *disgord.MessageCreate, args []string, client *disgord.Client) {

// 	defaultHelp(data, client)

// }

// NewHelpCommandClient returns a new instance of the HelpCommandClient
func NewHelpCommandClient(data *disgord.MessageCreate, disgordClient discordiface.DisgordClientAPI) *HelpCommandClient {
	return &HelpCommandClient{
		data:          data,
		disgordClient: disgordClient,
	}
}

// SendHelpMsg sends the default help message to the channel that received the help command
func (hcc *HelpCommandClient) SendHelpMsg() {
	resp, err := client.SendMsg(
		ctx,
		hcc.data.Message.ChannelID,
		&disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title: "**__help__**\n **ALIASES:** h, ?, wtf",
				Description: fmt.Sprintf(
					"This is the help function.\n\n"+
						"Use `%shelp functionName` to find out more about each function\n"+
						"Current available functions : ```\nhelp \n```"+
						"You can also read the source code here : https://github.com/aplombomb/boombot",
					conf.Prefix,
				),
				Footer: &disgord.EmbedFooter{
					Text: fmt.Sprintf("Help requested by %s", hcc.data.Message.Author.Username),
				},
				Timestamp: hcc.data.Message.Timestamp,
				Color:     0xeec400,
			},
		},
	)
	if err != nil {
		fmt.Println("There was an error sending default help message: ", err)
	}
	go deleteMessage(resp, 30*time.Second, client)
}
