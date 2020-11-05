package discord

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

type UnknownHandler struct {
	data          *disgord.Message
	disgordClient *disgord.Client
}

func NewUnknownHandler(data *disgord.Message, disgordClient *disgord.Client) *UnknownHandler {
	return &UnknownHandler{
		data:          data,
		disgordClient: disgordClient,
	}
}

func (uh *UnknownHandler) RespondToAuthor() error {

	resp, err := uh.disgordClient.SendMsg(
		ctx,
		uh.data.ChannelID,
		&disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:       "Unknown command",
				Description: fmt.Sprintf("Type %shelp to see the commands available", conf.Prefix),
				Timestamp:   uh.data.Timestamp,
				Color:       0xcc0000,
			},
		},
	)

	if err != nil {
		return err
	}
	// panic("\n\n\nMEEEEEP\n\n\n")
	go deleteMessage(uh.data, 150*time.Millisecond, client)
	go deleteMessage(resp, 10*time.Second, client)
	fmt.Println("Unknown command used")

	return nil
}
