package discord

import (
	"context"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
)

type DisgordClientAPI interface {
	CreateMessage(ctx context.Context, channelID snowflake.Snowflake, params *disgord.CreateMessageParams, flags ...disgord.Flag) (ret *disgord.Message, err error)
}
