package discord

import (
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
)

// MessageEventClient contains the data necessary for handling all non-command messages
type MessageEventClient struct {
	data          *disgord.Message
	disgordClient disgordiface.DisgordClientAPI
}

// NewMessageEventClient return a new MessageEventClient
func NewMessageEventClient(data *disgord.Message, disgordClient disgordiface.DisgordClientAPI) *MessageEventClient {
	return &MessageEventClient{
		data,
		disgordClient,
	}
}

//FilterNonModLinks removes all messages from mod requests channel that are not acceptable links for the minecraft mod suggestions channel
func (mec *MessageEventClient) FilterNonModLinks() error {
	//Per channel message event switch handler
	switch mec.data.ChannelID {
	case 734986357583380510:
		if !strings.Contains(mec.data.Content, "https://www.curseforge.com/minecraft/mc-mods/") {
			go deleteMessage(mec.data, 2*time.Second, mec.disgordClient)
		}
	case 851485354589814805:
		if strings.Contains(mec.data.Content, "https://external-preview.redd.it") {
			go deleteMessage(mec.data, 1*time.Second, mec.disgordClient)
		}
	default:
		break
	}
	return nil
}

func deleteMessage(resp *disgord.Message, sleep time.Duration, client disgordiface.DisgordClientAPI) {
	time.Sleep(sleep)
	// fmt.Printf("\nDeleting message '%+v' by user %+v", resp.Content, resp.Author.Username)
	channel := client.Channel(resp.ChannelID)
	msgQueryBuilder := channel.Message(resp.ID)
	msgQueryBuilder.Delete()
}

// SendEmbedMsgReply sends an embeded message
func (mec *MessageEventClient) SendEmbedMsgReply(embed disgord.Embed) (*disgord.Message, error) {
	resp, err := mec.disgordClient.SendMsg(
		mec.data.ChannelID,
		embed,
	)
	go deleteMessage(resp, 10*time.Second, mec.disgordClient)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
