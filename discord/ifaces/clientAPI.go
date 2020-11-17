package discord

import (
	"context"

	"github.com/andersfylling/disgord"
)

// DisgordClientAPI provides an interface for mocking disgord client behavior
type DisgordClientAPI interface {
	SendMsg(channelID disgord.Snowflake, data ...interface{}) (msg *disgord.Message, err error)
	GetMessage(ctx context.Context, channelID, messageID disgord.Snowflake, flags ...disgord.Flag) (message *disgord.Message, err error)
	VoiceConnect(guildID, channelID disgord.Snowflake) (disgord.VoiceConnection, error)
	Delete(flags ...disgord.Flag) (err error)
}

type DisgordMessageAPI interface {
}

// DisgordUserAPI provides an interface for mocking disgord user behavior
type DisgordUserAPI interface {
	SendMsg(ctx context.Context, session disgord.Session, message *disgord.Message) (channel *disgord.Channel, msg *disgord.Message, err error)
}
