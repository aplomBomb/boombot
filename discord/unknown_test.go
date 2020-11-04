package discord_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/andersfylling/disgord"
	newclientmock "github.com/aplombomb/boombot/_mocks/generated/discord/newclient"
	sendmsgmock "github.com/aplombomb/boombot/_mocks/generated/discord/sendmsg"

	"github.com/aplombomb/boombot/discord"
)

func TestUnknown(t *testing.T) {

	c := gomock.NewController(t)
	mockNewClient := newclientmock.NewMockDisgordClientAPI(c)
	mockNewClient.EXPECT().NewClient(gomock.Any()).Return(&disgord.Client{}, nil)
	mockedClient, _ := mockNewClient.NewClient(disgord.Config{})

	fmt.Printf("\n\n\nMOCKEDCLIENT: %+v\n\n\n,", mockedClient)

	type args struct {
		data   *disgord.Message
		client *disgord.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "unknown | send msg success",
			args: args{
				&disgord.Message{
					ChannelID: 4234234234,
					Timestamp: disgord.Time{Time: time.Now()},
					ID:        2342342343,
				},
				mockedClient,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			mockNewClient := newclientmock.NewMockDisgordClientAPI(c)
			mockNewClient.EXPECT().NewClient(gomock.Any()).Return(&disgord.Client{}, nil)
			mockedClient, _ := mockNewClient.NewClient(disgord.Config{})
			mockSendMsg := sendmsgmock.NewMockDisgordMsgAPI(c)

			mockSendMsg.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(&disgord.Message{}, nil)

			if err := discord.Unknown(tt.args.data, mockedClient); (err != nil) != tt.wantErr {
				t.Errorf("Unknown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
