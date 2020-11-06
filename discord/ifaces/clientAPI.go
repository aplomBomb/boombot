package discord

import (
	"context"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
)

// DisgordClientAPI provides an interface for mocking disgord client behavior
type DisgordClientAPI interface {
	SendMsg(ctx context.Context, channelID snowflake.Snowflake, data ...interface{}) (msg *disgord.Message, err error)
	GetMessage(ctx context.Context, channelID, messageID snowflake.Snowflake, flags ...disgord.Flag) (message *disgord.Message, err error)
	DeleteMessage(ctx context.Context, channelID, msgID snowflake.Snowflake, flags ...disgord.Flag) (err error)
}
