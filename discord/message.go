package discord

import (
	"strings"

	"github.com/andersfylling/disgord"
)

//RespondToMessage handles all messages created in the server
func RespondToMessage(s disgord.Session, data *disgord.MessageCreate) {
	//Per channel message event switch handler
	switch data.Message.ChannelID {
	case 734986357583380510:
		if strings.Contains(data.Message.Content, "https://www.curseforge.com/minecraft/mc-mods/") == false {
			message, _ := client.GetMessage(ctx, data.Message.ChannelID, data.Message.ID)
			go deleteMessage(message, 1)
		}
	default:
		break
	}

}

// ParseMessage parses the message into command / args
func ParseMessage(data *disgord.MessageCreate) (string, []string) {
	var command string
	var args []string
	if len(data.Message.Content) > 0 {
		command = strings.ToLower(strings.Fields(data.Message.Content)[0])
		if len(data.Message.Content) > 1 {
			args = strings.Fields(data.Message.Content)[1:]
		}
	}
	return command, args
}
