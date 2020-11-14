package discord

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

// TO-DO All of the voice related logic will be moved to the yt package
// Getting everything working in here first

// VoiceChannel defines a voice channel's current user state and ID
type VoiceChannel struct {
	ID    disgord.Snowflake
	Name  string
	Users []*disgord.User
}

// VoiceChannels contains a collection of VoiceChannel
type VoiceChannels struct {
	Channels []*VoiceChannel
}

var voiceChannelCache VoiceChannels

// Using this for access to the global clients FOR NOW as passing it through the handlers has proven tricky
// TO-DO find a solution to get rid of the global variables, including the client

// RespondToCommand delegates actions when commands are issued
func RespondToCommand(s disgord.Session, data *disgord.MessageCreate) {
	fmt.Printf("\nvoiceChannelCache: %+v\n", voiceChannelCache)
	for i, v := range voiceChannelCache.Channels {
		fmt.Printf("\nvoiceChannelCache channel %+v users: %+v\n", i, v.Users)
	}

	user, err := disgordGlobalClient.GetUser(ctx, data.Message.Author.ID)
	if err != nil {
		fmt.Println("Failed to fetch user (probably a webhook)")
		user = &disgord.User{
			Username: "unknown",
		}
	}
	cec := NewCommandEventClient(data.Message, disgordGlobalClient, ytService.Search)
	command, _ := cec.DisectCommand()
	fmt.Printf("Command %+v by user %+v | %+v\n", command, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	cec.Delegate()
}

// RespondToMessage delegates actions when messages are created
func RespondToMessage(s disgord.Session, data *disgord.MessageCreate) {
	user, err := disgordGlobalClient.GetUser(ctx, data.Message.Author.ID)
	if err != nil {
		fmt.Println("Failed to fetch user (probably a webhook)")
		user = &disgord.User{
			Username: "unknown",
		}
	}
	fmt.Printf("Message %+v by user %+v | %+v\n", data.Message.Content, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	mec := NewMessageEventClient(data.Message, disgordGlobalClient)
	err = mec.FilterNonModLinks()
	if err != nil {
		fmt.Printf("\nError filtering non-mod link: %+v\n", err)
	}
}

// RespondToReaction delegates actions when reactions are added to messages
func RespondToReaction(s disgord.Session, data *disgord.MessageReactionAdd) {
	user, _ := disgordGlobalClient.GetUser(ctx, data.UserID)
	// fmt.Printf("Message reaction %+v by user %+v | %+v\n", data.PartialEmoji.Name, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	rec := NewReactionEventClient(data.PartialEmoji, data.UserID, data.ChannelID, data.MessageID, disgordGlobalClient)
	msg, err := rec.GenerateModResponse()
	if err != nil {
		fmt.Printf("\nError generating mod reaction response: %+v\n", err)
	}
	//TO-DO Sending the dm here as opposed having it sent via GenerateModResponse for testing purposes
	// Using it here at least allows me to get full coverage of the reactions logic
	// The SendMsg method of disgord.User requires session arg which has proven difficult to mock
	if msg != nil {
		user.SendMsg(ctx, s, msg)
	}
}

// RespondToVoiceChannelUpdate updates the server's voice channel member cache every time an update is emitted
func RespondToVoiceChannelUpdate(s disgord.Session, data *disgord.VoiceStateUpdate) {
	// TO-DO make cache into a map: [userID]channelID
	// Delete entry when user leaves a voice channel \ when channelID on event is 0
	newVoiceChannelCache := VoiceChannels{}
	channels, err := s.GetGuildChannels(ctx, data.GuildID)
	if err != nil {
		fmt.Printf("\nError getting guild channels: %+v\n", err)
	}
	for _, v := range channels {
		if v.Type == 2 {
			channel, err := s.GetChannel(ctx, v.ID)
			if err != nil {
				fmt.Printf("\nError getting guild channel: %+v\n", err)
			}
			fmt.Printf("\nVoice Channel recipients: %+v", channel.Recipients)
			newChannelDetails := &VoiceChannel{
				ID:    v.ID,
				Name:  v.Name,
				Users: channel.Recipients,
			}
			newVoiceChannelCache.Channels = append(newVoiceChannelCache.Channels, newChannelDetails)
			voiceChannelCache = newVoiceChannelCache
		}
	}
	u, err := data.Member.GetUser(ctx, s)
	if err != nil {
		fmt.Printf("\nError getting user: %+v\n", err)
	}
	fmt.Printf("\nUserObject: %+v", u)

	// s.VoiceConnect(data.GuildID, data.ChannelID)

}
