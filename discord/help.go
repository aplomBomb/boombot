package discord

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
)

// HelpCommandClient contains the resources needed for handling help requests
type HelpCommandClient struct {
	data          *disgord.Message
	disgordClient disgordiface.DisgordClientAPI
}

// NewHelpCommandClient returns a new instance of the HelpCommandClient
func NewHelpCommandClient(data *disgord.Message, disgordClient disgordiface.DisgordClientAPI) *HelpCommandClient {
	return &HelpCommandClient{
		data:          data,
		disgordClient: disgordClient,
	}
}

// SendHelpMsg sends the default help message to the channel that received the help command
func (hcc *HelpCommandClient) SendHelpMsg() error {
	resp, err := hcc.disgordClient.SendMsg(
		// ctx,
		hcc.data.ChannelID,
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
					Text: fmt.Sprintf("Help requested by %s", hcc.data.Author.Username),
				},
				Timestamp: hcc.data.Timestamp,
				Color:     0xeec400,
			},
		},
	)
	if err != nil {
		return err
	}
	go deleteMessage(hcc.data, 150*time.Millisecond, hcc.disgordClient)
	go deleteMessage(resp, 30*time.Second, hcc.disgordClient)
	return nil
}
