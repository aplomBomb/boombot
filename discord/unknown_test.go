package discord_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/andersfylling/disgord"
	mock_sendmsg "github.com/aplombomb/boombot/_mocks/generated/discord"
	"github.com/aplombomb/boombot/discord"
)

func TestUnknown(t *testing.T) {
	type args struct {
		data *disgord.Message
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			mockSendMsg := mock_sendmsg.NewMockDisgordMsgAPI(c)
			// mockSendMsg.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(&disgord.Message{}, nil)

			if err := discord.Unknown(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Unknown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
