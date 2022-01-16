package discord_test

import (
	"testing"

	"github.com/andersfylling/disgord"
	clientMock "github.com/aplombomb/boombot/_mocks/generated/discordclient"
	"github.com/aplombomb/boombot/discord"
	discordiface "github.com/aplombomb/boombot/discord/ifaces"
	"github.com/golang/mock/gomock"
)

func TestMessageEventClient_FilterNonModLinks(t *testing.T) {
	c := gomock.NewController(t)
	mc := clientMock.NewMockDisgordClientAPI(c)

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
			name: "message | non valid link",
			fields: func() fields {

				return fields{
					data: &disgord.Message{
						Content:   "wrong stuff!",
						ChannelID: 734986357583380510,
					},
					disgordClient: mc,
				}
			}(),
			wantErr: false,
		},
		{
			name: "message | valid link",
			fields: func() fields {

				return fields{
					data: &disgord.Message{
						Content:   "https://www.curseforge.com/minecraft/mc-mods/",
						ChannelID: 734986357583380510,
					},
					disgordClient: mc,
				}
			}(),
			wantErr: false,
		},
		{
			name: "message | invalid channel",
			fields: func() fields {

				return fields{
					data: &disgord.Message{
						Content:   "https://www.curseforge.com/minecraft/mc-mods/",
						ChannelID: 7,
					},
					disgordClient: mc,
				}
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mec := discord.NewMessageEventClient(tt.fields.data, tt.fields.disgordClient)
			if err := mec.FilterMessages(); (err != nil) != tt.wantErr {
				t.Errorf("MessageEventClient.FilterNonModLinks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
