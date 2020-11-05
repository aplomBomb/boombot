package discord_test

import (
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/andersfylling/disgord"
	mockmsg "github.com/aplombomb/boombot/_mocks/generated/discord/sendmsg"
	"github.com/aplombomb/boombot/discord"
)

// func defaultFields(t *testing.T) fields {
// 	c := gomock.NewController(t)
// 	clientMock := mock
// 	return fields{}
// }

func TestUnknownHandler_RespondToAuthor(t *testing.T) {

	message := &disgord.Message{
		ChannelID: 645286762452877314,
		Timestamp: disgord.Time{Time: time.Now()},
		Content:   "**THIS IS A TEST FOR THE UNKOWN COMMAND**",
	}

	type fields struct {
		data          *disgord.Message
		disgordClient *disgord.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "RespondToAuthor | success",
			fields: func() fields {
				c := gomock.NewController(t)
				mc := mockmsg.NewMockDisgordMsgAPI(c)
				mc.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(&disgord.Message{}, nil)
				return fields{
					data:          message,
					disgordClient: disgord.New(disgord.Config{BotToken: os.Getenv("BOOMBOT_TOKEN")}),
				}
			}(),
			wantErr: false,
		},
		// {
		// 	name: "RespondToAuthor | SendMsg error",
		// 	fields: func() fields {
		// 		c := gomock.NewController(t)
		// 		mc := mockmsg.NewMockDisgordMsgAPI(c)
		// 	}
		// }
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uh := discord.NewUnknownHandler(tt.fields.data, tt.fields.disgordClient)
			if err := uh.RespondToAuthor(); (err != nil) != tt.wantErr {
				t.Errorf("UnknownHandler.RespondToAuthor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
