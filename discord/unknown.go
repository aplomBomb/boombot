package discord

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

func unknown(data *disgord.MessageCreate, message *disgord.Message) *disgord.Message {
	go deleteMessage(data.Message, 1)
	go deleteMessage(message, 10)
	fmt.Println("Unknown command used")

	return message
}
