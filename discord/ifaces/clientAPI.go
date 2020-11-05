package discord

import (
	"context"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
)

type DisgordClientAPI interface {
	SendMsg(ctx context.Context, channelID snowflake.Snowflake, data ...interface{}) (msg *disgord.Message, err error)
}
