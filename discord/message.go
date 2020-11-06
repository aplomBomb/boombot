package discord

import (
	"strings"

	"github.com/andersfylling/disgord"
	discord "github.com/aplombomb/boombot/discord/ifaces"
)

type MessageEventClient struct {
	sess          *disgord.Session
	data          *disgord.Message
	disgordClient discord.DisgordClientAPI
}

func NewMessageEventClient(sess *disgord.Session, data *disgord.Message, disgordClient discord.DisgordClientAPI) *MessageEventClient {
	return &MessageEventClient{
		sess,
		data,
		disgordClient,
	}
}

//RespondToMessage handles all messages created in the server
func (mec *MessageEventClient) RespondToMessage() {
	//Per channel message event switch handler
	switch mec.data.ChannelID {
	case 734986357583380510:
		if strings.Contains(mec.data.Content, "https://www.curseforge.com/minecraft/mc-mods/") == false {
			message, _ := mec.disgordClient.GetMessage(ctx, mec.data.ChannelID, mec.data.ID)
			go deleteMessage(message, 1, mec.disgordClient)
		}
	default:
		break
	}

}

// ParseMessage parses the message into command / args
func ParseMessage(data *disgord.Message) (string, []string) {
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
