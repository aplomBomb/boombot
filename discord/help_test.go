package discord_test

// import (
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/golang/mock/gomock"

// 	"github.com/andersfylling/disgord"

// 	mockClient "github.com/aplombomb/boombot/_mocks/generated/discordclient"
// 	"github.com/aplombomb/boombot/discord"
// 	discordiface "github.com/aplombomb/boombot/discord/ifaces"
// )

// func TestHelpCommandClient_SendHelpMsg(t *testing.T) {
// 	c := gomock.NewController(t)
// 	mc := mockClient.NewMockDisgordClientAPI(c)
// 	testMessage := &disgord.Message{
// 		ChannelID: 123,
// 		Timestamp: disgord.Time{Time: time.Now()},
// 		Author: &disgord.User{
// 			Username: "bomb",
// 		},
// 	}

// 	type fields struct {
// 		data          *disgord.Message
// 		disgordClient discordiface.DisgordClientAPI
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		{
// 			name: "help | sendmsg success",
// 			fields: func() fields {
// 				mc.EXPECT().SendMsg(gomock.Any(), gomock.Any(), gomock.Any()).Return(&disgord.Message{}, nil)
// 				return fields{
// 					testMessage,
// 					mc,
// 				}
// 			}(),
// 			wantErr: false,
// 		},
// 		{
// 			name: "help | sendmsg error",
// 			fields: func() fields {
// 				mc.EXPECT().SendMsg(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("help sndmsgERR"))
// 				return fields{
// 					testMessage,
// 					mc,
// 				}
// 			}(),
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			hcc := discord.NewHelpCommandClient(tt.fields.data, tt.fields.disgordClient)

// 			if err := hcc.SendHelpMsg(); (err != nil) != tt.wantErr {
// 				t.Errorf("HelpCommandClient.SendHelpMsg() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
