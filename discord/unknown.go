package discord

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

type UnknownCommandClient struct {
	data          *disgord.Message
	disgordClient *disgord.Client
}

func NewUnknownCommandClient(data *disgord.Message, disgordClient *disgord.Client) *UnknownCommandClient {
	return &UnknownCommandClient{
		data:          data,
		disgordClient: disgordClient,
	}
}

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
