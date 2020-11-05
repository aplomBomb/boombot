package discord

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
	discord "github.com/aplombomb/boombot/discord/ifaces"
)

// UnknownCommandClient represents the data neccessary for unknown command processing
type UnknownCommandClient struct {
	data          *disgord.Message
	disgordClient discord.DisgordClientAPI
}

// NewUnknownCommandClient returns a new instance
func NewUnknownCommandClient(data *disgord.Message, disgordClient discord.DisgordClientAPI) *UnknownCommandClient {
	return &UnknownCommandClient{
		data:          data,
		disgordClient: disgordClient,
	}
}

// RespondToChannel handles sending a message to the channel that received an unknown command
func (uc *UnknownCommandClient) RespondToChannel() error {

	resp, err := uc.disgordClient.SendMsg(
		ctx,
		uc.data.ChannelID,
		&disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:       "Unknown command",
				Description: fmt.Sprintf("Type %shelp to see the commands available", conf.Prefix),
				Timestamp:   uc.data.Timestamp,
				Color:       0xcc0000,
			},
		},
	)

	if err != nil {
		return err
	}
	// panic("\n\n\nMEEEEEP\n\n\n")
	go deleteMessage(uc.data, 150*time.Millisecond, client)
	go deleteMessage(resp, 10*time.Second, client)
	fmt.Println("Unknown command used")

	return nil
}
