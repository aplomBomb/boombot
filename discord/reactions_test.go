package discord_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/andersfylling/disgord"
	clientMock "github.com/aplombomb/boombot/_mocks/generated/discordclient"
	"github.com/aplombomb/boombot/discord"
	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
)

func TestReactionEventClient_GenerateModResponse(t *testing.T) {
	c := gomock.NewController(t)
	mockClient := clientMock.NewMockDisgordClientAPI(c)

	seenMessage := &disgord.Message{
		Embeds: []*disgord.Embed{
			&disgord.Embed{
				Title:       fmt.Sprintf("**Your request to add %s is being reviewed**", "testMod"),
				URL:         "https://www.curseforge.com/minecraft/mc-mods/testMod",
				Description: fmt.Sprintf("*Bomb is reviewing your request to add %s*", "testMod"),
				Color:       0xcc0000,
				Footer: &disgord.EmbedFooter{
					Text:    "Sit tight partner!",
					IconURL: "https://cdn.discordapp.com/emojis/745396324215685201.gif?v=1",
				},
			},
		},
	}

	acceptedMessage := &disgord.Message{
		Embeds: []*disgord.Embed{
			&disgord.Embed{
				Title:       fmt.Sprintf("**%s ACCEPTED!!**", "testMod"),
				URL:         "https://www.curseforge.com/minecraft/mc-mods/testMod",
				Description: fmt.Sprintf("*Bomb has added %s to the modpack! If the server breaks now, it's all your fault!*", "testMod"),
				Color:       0xcc0000,
				Footer: &disgord.EmbedFooter{
					Text:    "Pervert Steve is always watching...",
					IconURL: "https://cdn.discordapp.com/emojis/681217726412488767.png?v=1",
				},
			},
		},
	}

	rejectedMessage := &disgord.Message{
		Embeds: []*disgord.Embed{
			&disgord.Embed{
				Title:       fmt.Sprintf("**%s Rejected**", "testMod"),
				URL:         "https://www.curseforge.com/minecraft/mc-mods/testMod",
				Description: fmt.Sprintf("*Bomb has rejected your request to add %s*", "testMod"),
				Color:       0xcc0000,
				Footer: &disgord.EmbedFooter{
					Text:    "You have brought much shame upon your famiry",
					IconURL: "https://cdn.discordapp.com/emojis/662170922580574258.gif?v=1",
				},
			},
		},
	}

	type fields struct {
		emoji         *disgord.Emoji
		uID           disgord.Snowflake
		chID          disgord.Snowflake
		msgID         disgord.Snowflake
		disgordClient disgordiface.DisgordClientAPI
	}
	tests := []struct {
		name    string
		fields  fields
		want    *disgord.Message
		wantErr bool
	}{
		{
			name: "seen reaction | make message success",
			fields: func() fields {
				mockClient.EXPECT().GetMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(&disgord.Message{
					Content: "https://www.curseforge.com/minecraft/mc-mods/testMod",
				}, nil)
				return fields{
					emoji: &disgord.Emoji{
						Name: "eyes",
					},
					uID:           321044596476084235,
					chID:          734986357583380510,
					msgID:         456,
					disgordClient: mockClient,
				}
			}(),
			want:    seenMessage,
			wantErr: false,
		},
		{
			name: "seen reaction | get message error",
			fields: func() fields {
				mockClient.EXPECT().GetMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("Seen reaction | GetMessage error"))
				return fields{
					emoji: &disgord.Emoji{
						Name: "eyes",
					},
					uID:           321044596476084235,
					chID:          734986357583380510,
					msgID:         456,
					disgordClient: mockClient,
				}
			}(),
			want:    nil,
			wantErr: true,
		},
		{
			name: "accepted reaction | make message success",
			fields: func() fields {
				mockClient.EXPECT().GetMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(&disgord.Message{
					Content: "https://www.curseforge.com/minecraft/mc-mods/testMod",
				}, nil)
				return fields{
					emoji: &disgord.Emoji{
						Name: "check",
					},
					uID:           321044596476084235,
					chID:          734986357583380510,
					msgID:         456,
					disgordClient: mockClient,
				}
			}(),
			want:    acceptedMessage,
			wantErr: false,
		},
		{
			name: "accepted reaction | get message error",
			fields: func() fields {
				mockClient.EXPECT().GetMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("Accepted reaction | GetMessage error"))
				return fields{
					emoji: &disgord.Emoji{
						Name: "check",
					},
					uID:           321044596476084235,
					chID:          734986357583380510,
					msgID:         456,
					disgordClient: mockClient,
				}
			}(),
			want:    nil,
			wantErr: true,
		},
		{
			name: "rejected reaction | make message success",
			fields: func() fields {
				mockClient.EXPECT().GetMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(&disgord.Message{
					Content: "https://www.curseforge.com/minecraft/mc-mods/testMod",
				}, nil)
				return fields{
					emoji: &disgord.Emoji{
						Name: "no",
					},
					uID:           321044596476084235,
					chID:          734986357583380510,
					msgID:         456,
					disgordClient: mockClient,
				}
			}(),
			want:    rejectedMessage,
			wantErr: false,
		},
		{
			name: "rejected reaction | get message error",
			fields: func() fields {
				mockClient.EXPECT().GetMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("Rejected reaction | GetMessage error"))
				return fields{
					emoji: &disgord.Emoji{
						Name: "no",
					},
					uID:           321044596476084235,
					chID:          734986357583380510,
					msgID:         456,
					disgordClient: mockClient,
				}
			}(),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := discord.NewReactionEventClient(
				tt.fields.emoji,
				tt.fields.uID,
				tt.fields.chID,
				tt.fields.msgID,
				tt.fields.disgordClient,
			)
			got, err := rec.GenerateModResponse()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReactionEventClient.GenerateModResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReactionEventClient.GenerateModResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
