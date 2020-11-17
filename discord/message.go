package discord

import (
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

// MessageEventClient contains the data necessary for handling all non-command messages
type MessageEventClient struct {
	data          *disgord.Message
	disgordClient *disgord.Client
}

// NewMessageEventClient return a new MessageEventClient
func NewMessageEventClient(data *disgord.Message, disgordClient *disgord.Client) *MessageEventClient {
	return &MessageEventClient{
		data,
		disgordClient,
	}
}

//FilterNonModLinks removes all messages from mod requests channel that are not acceptable links
func (mec *MessageEventClient) FilterNonModLinks() error {
	//Per channel message event switch handler
	switch mec.data.ChannelID {
	case 734986357583380510:
		if strings.Contains(mec.data.Content, "https://www.curseforge.com/minecraft/mc-mods/") == false {
			go deleteMessage(mec.data, 2*time.Second, mec.disgordClient)
		}
	default:
		break
	}
	return nil
}
