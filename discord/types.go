package discord

import (
	"context"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
)

type DisgordMsgAPI interface {
	SendMsg(ctx context.Context, channelID snowflake.Snowflake, data ...interface{}) (msg *disgord.Message, err error)
}

type DisgordClientAPI interface {
	NewClient(conf disgord.Config) (*disgord.Client, error)
}
