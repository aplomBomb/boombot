package discord

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

func unknown(data *disgord.Message) (*disgord.Message, error) {

	resp, err := client.SendMsg(
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
		return nil, err
	}

	fmt.Printf("\nTHIS IS THE MESSAGE: %+v\n", resp.Content)

	go deleteMessage(data, 150*time.Millisecond)
	go deleteMessage(resp, 10*time.Second)
	fmt.Println("Unknown command used")

	return resp, nil
}
