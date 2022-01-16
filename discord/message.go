package discord

import (
	"fmt"
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

//FilterMessages is used for deleting unwanted messages from the channel of origin
func (mec *MessageEventClient) FilterMessages() error {
	//Per channel message event switch handler
	switch mec.data.ChannelID {
	case ServerIDs.McModChID:
		if !strings.Contains(mec.data.Content, "https://www.curseforge.com/minecraft/mc-mods/") {
			go deleteMessage(mec.data, 2*time.Second, mec.disgordClient)
		}
	case ServerIDs.TihiID:
		if strings.Contains(mec.data.Content, "https://external-preview.redd.it") {
			fmt.Println("\nBINGO!")
			go deleteMessage(mec.data, 1*time.Second, mec.disgordClient)
		}
	default:
		break
	}
	return nil
}

func deleteMessage(resp *disgord.Message, sleep time.Duration, client disgordiface.DisgordClientAPI) {
	time.Sleep(sleep)
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
