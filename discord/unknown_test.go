package discord_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/andersfylling/disgord"
	mockSendMsg "github.com/aplombomb/boombot/_mocks/generated/discordclient"
	"github.com/aplombomb/boombot/discord"
	discordiface "github.com/aplombomb/boombot/discord/ifaces"
)

func TestUnknownCommandClient_RespondToChannel(t *testing.T) {
	c := gomock.NewController(t)
	msm := mockSendMsg.NewMockDisgordClientAPI(c)

	testMessage := &disgord.Message{
		ChannelID: 123,
		Timestamp: disgord.Time{Time: time.Now()},
	}

	type fields struct {
		data          *disgord.Message
		disgordClient discordiface.DisgordClientAPI
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "unknowncommandclient | respondToChannel success",
			fields: func() fields {
				msm.EXPECT().SendMsg(gomock.Any(), gomock.Any(), gomock.Any()).Return(&disgord.Message{}, nil)
				return fields{
					testMessage,
					msm,
				}
			}(),
			wantErr: false,
		},
		{
			name: "unknowncommandclient | respondToChannel error",
			fields: func() fields {
				msm.EXPECT().SendMsg(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("unknownsendmsgERR"))
				return fields{
					testMessage,
					msm,
				}
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := discord.NewUnknownCommandClient(tt.fields.data, tt.fields.disgordClient)

			if err := uc.RespondToChannel(); (err != nil) != tt.wantErr {
				t.Errorf("UnknownCommandClient.RespondToChannel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
