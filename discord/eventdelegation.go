package discord

import (
	"fmt"
	"os"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/jonas747/dca"
)

// TO-DO All of the voice related logic will be moved to the yt package
// Getting everything working in here first
var voiceChannelCache = make(map[disgord.Snowflake]disgord.Snowflake)

// Using this for access to the global clients FOR NOW as passing it through the handlers has proven tricky
// TO-DO find a solution to get rid of the global variables, including the client

// RespondToCommand delegates actions when commands are issued
func RespondToCommand(s disgord.Session, data *disgord.MessageCreate) {
	cec := NewCommandEventClient(data.Message, disgordGlobalClient)
	command, _ := cec.DisectCommand()

	// fmt.Printf("\nvCommand: %+v | Args: %+v\n", command, args)

	user, err := disgordGlobalClient.GetUser(ctx, data.Message.Author.ID)
	if err != nil {
		fmt.Println("Failed to fetch user (probably a webhook)")
		user = &disgord.User{
			Username: "unknown",
		}
	}

	fmt.Printf("Command %+v by user %+v | %+v\n", command, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	switch command {
	case "play":

		// ss := youtube.NewSearchService(ytService)
		// ytc := yt.NewYoutubeClient(ss)
		// payload, _ := ytc.SearchAndDownload(args[0])

		// out, err := os.Create(fmt.Sprintf("%+v.mp3", args[0]))
		// if err != nil {
		// 	// panic?
		// }

		encodeSess, err := dca.EncodeFile("tKi9Z-f6qX4.mp3", dca.StdEncodeOptions)
		if err != nil {
			fmt.Printf("\nERROR ENCODING: %+v\n", err)
		}

		// io.Copy(out, payload.Body)

		vc, err := s.VoiceConnect(data.Message.GuildID, 737468810222895125)
		if err != nil {
			fmt.Printf("\nERROR CONNECTING TO VOICE CHANNEL: %+v\n", err)
		}
		err = vc.StartSpeaking()
		if err != nil {
			fmt.Printf("\nERROR SPEAKING: %+v\n", err)
		}
		song, err := os.Open("./discord/strobe.dca")
		fmt.Printf("\nERROR OPENING: %+v\n", err)
		fmt.Printf("\nSONG: %+v\n", song)
		err = vc.SendDCA(encodeSess)
		if err != nil {
			fmt.Printf("\nERROR PLAYING DCA: %+v\n", err)
		}

		time.Sleep(5 * time.Second)
		vc.StopSpeaking()
		time.Sleep(5 * time.Second)
		s.Disconnect()

		// defer out.Close()
		// defer encodeSess.Cleanup()
		// defer payload.Body.Close()

	default:
		cec.Delegate()
	}

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
	channel, err := s.GetChannel(ctx, data.ChannelID)
	if err != nil {
		fmt.Printf("\nError getting channel: %+v\n", err)
	}

	switch data.ChannelID {
	// If ChannelID is 0, then a user left a channel, delete them from the cache
	case 0:
		fmt.Printf("\nUser left %+v\n", voiceChannelCache[data.UserID])
		delete(voiceChannelCache, data.UserID)
		// Add userID and ChannelID to the voiceChannelCache upon join
	default:
		fmt.Printf("\nUser joined %+v\n", channel.Name)
		voiceChannelCache[data.UserID] = channel.ID
	}

}
