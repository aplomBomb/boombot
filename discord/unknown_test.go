package discord_test

import (
	"os"
	"testing"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/aplombomb/boombot/discord"
)

func TestUnknownHandler_RespondToAuthor(t *testing.T) {

	testMessage := &disgord.Message{
		ChannelID: 789789789,
		Timestamp: disgord.Time{Time: time.Now()},
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
			name: "RespondToAuthor | channelID error",
			fields: func() fields {

				return fields{
					data:          testMessage,
					disgordClient: disgord.New(disgord.Config{BotToken: os.Getenv("BOOMBOT_TOKEN")}),
				}
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uh := discord.NewUnknownCommandClient(tt.fields.data, tt.fields.disgordClient)
			if err := uh.RespondToChannel(); (err != nil) != tt.wantErr {
				t.Errorf("UnknownHandler.RespondToAuthor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
