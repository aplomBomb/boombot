package discord

import (
	"strings"

	"github.com/andersfylling/disgord"
	discord "github.com/aplombomb/boombot/discord/ifaces"
)

// MessageEventClient contains the data necessary for handling all non-command messages
type MessageEventClient struct {
	data          *disgord.Message
	disgordClient discord.DisgordClientAPI
}

// NewMessageEventClient return a new MessageEventClient
func NewMessageEventClient(data *disgord.Message, disgordClient discord.DisgordClientAPI) *MessageEventClient {
	return &MessageEventClient{
		data,
		disgordClient,
	}
}

//FilterNonModLinks removes all messages from mod requests channel that are not acceptable links
func (mec *MessageEventClient) FilterNonModLinks() {
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
