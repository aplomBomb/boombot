package discord

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

func Unknown(data *disgord.Message, disgordClient *disgord.Client) error {

	resp, err := disgordClient.SendMsg(
		ctx,
		data.ChannelID,
		&disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title:       "Unknown command",
				Description: fmt.Sprintf("Type %shelp to see the commands available", conf.Prefix),
				Timestamp:   data.Timestamp,
				Color:       0xcc0000,
			},
		},
	)

	if err != nil {
		return err
	}
	// panic("\n\n\nMEEEEEP\n\n\n")
	go deleteMessage(data, 150*time.Millisecond, client)
	go deleteMessage(resp, 10*time.Second, client)
	fmt.Println("Unknown command used")

	return nil
}
