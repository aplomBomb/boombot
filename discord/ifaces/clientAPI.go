package discord

import (
	"context"

	"github.com/andersfylling/disgord"
)

// DisgordClientAPI provides an interface for mocking Disgord's cache behavior
type DisgordClientAPI interface {
	SendMsg(channelID disgord.Snowflake, data ...interface{}) (msg *disgord.Message, err error)
	// GetMessage(channelID, messageID disgord.Snowflake) (ret *disgord.Message, err error)
	VoiceConnectOptions(guildID, channelID disgord.Snowflake, selfDeaf, selfMute bool) (ret disgord.VoiceConnection, err error)
	Channel(id disgord.Snowflake) disgord.ChannelQueryBuilder
	User(id disgord.Snowflake) disgord.UserQueryBuilder
}

// DisgordMessageQueryBuilderAPI provides an interface for mocking Disgord's MessageQueryBuilder behavior
type DisgordMessageQueryBuilderAPI interface {
	Delete(flags ...disgord.Flag) (err error)
}

// DisgordChannelQueryBuilderAPI provides an interface for mocking Disgord's MessageQueryBuilder behavior
type DisgordChannelQueryBuilderAPI interface {
	Delete(flags ...disgord.Flag) (err error)
}

// DisgordUserAPI provides an interface for mocking disgord user behavior
type DisgordUserAPI interface {
	SendMsg(ctx context.Context, session disgord.Session, message *disgord.Message) (channel *disgord.Channel, msg *disgord.Message, err error)
}

// QueueClientAPI provides an interface for mocking boombot queue behavior
type QueueClientAPI interface {
	TriggerNext()
	TriggerShuffle()
	TriggerStop()
	UpdateUserQueueStateBulk(chID disgord.Snowflake, uID disgord.Snowflake, args []string)
	UpdateUserQueueState(chID disgord.Snowflake, uID disgord.Snowflake, arg string)
	TriggerChannelHop(id disgord.Snowflake)
	ReturnUserQueue() map[disgord.Snowflake][]string
	ReturnNowPlayingID() disgord.Snowflake
	ReturnVoiceCacheEntry(id disgord.Snowflake) disgord.Snowflake
}
