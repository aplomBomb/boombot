package discord

import (
	"github.com/andersfylling/disgord"
)

//RespondToVoiceChannelUpdate contains logic for handling the voiceChannelUpdate event
func RespondToVoiceChannelUpdate(s disgord.Session, data *disgord.VoiceStateUpdate) {

	// gateway, err := client.GetGateway(data.Ctx)

	// defer gateway.

	// if err != nil {
	// 	log.Fatal("\nError creating gateway: ", err)
	// }

	// var voice disgord.VoiceConnection

	// fmt.Printf("User %+v just joined the %+v voice chat", data.Member.User.Username, data.ChannelID)

	// fmt.Printf("\n\nctx: %+v | guildID: %+v \n\n", data.Ctx, data.GuildID)

	// channelSnowflakes, err := session.GetGuildChannels(data.Ctx, data.GuildID)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// client.VoiceConnect(data.GuildID, channelSnowflakes[0].ID)

}
