package discord

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/andersfylling/disgord"
	mockClient "github.com/aplombomb/boombot/_mocks/generated/discordclient"
	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
)

func TestReactionEventClient_RespondToReaction(t *testing.T) {
	c := gomock.NewController(t)
	clientMock := mockClient.NewMockDisgordClientAPI(c)
	type fields struct {
		emoji          *disgord.Emoji
		uID            disgord.Snowflake
		chID           disgord.Snowflake
		msgID          disgord.Snowflake
		disgordClient  disgordiface.DisgordClientAPI
		disgordSession disgord.Session
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "mod accepted success",
			fields: func() fields {
				return fields{
					emoji: &disgord.Emoji{
						Name: "testEmoji",
					},
					uID:            123,
					chID:           321,
					msgID:          456,
					disgordClient:  clientMock,
					disgordSession: session,
				}
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := NewReactionEventClient(
				tt.fields.emoji,
				tt.fields.uID,
				tt.fields.chID,
				tt.fields.msgID,
				tt.fields.disgordClient,
				tt.fields.disgordSession,
			)
			if err := rec.RespondToReaction(); (err != nil) != tt.wantErr {
				t.Errorf("ReactionEventClient.RespondToReaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
