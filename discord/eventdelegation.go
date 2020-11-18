package discord

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
	yt "github.com/aplombomb/boombot/youtube"
	"github.com/jonas747/dca"
	"google.golang.org/api/youtube/v3"
)

// TO-DO All of the voice related logic will be moved to the yt package
// Getting everything working in here first

// VoiceChannelCache keeps a map record of userID's and the channelID that user currently belongs to
type VoiceChannelCache struct {
	Cache map[disgord.Snowflake]disgord.Snowflake
}

// NewVoiceChannelCache return a pointer to a new voicechannelcache
func NewVoiceChannelCache() *VoiceChannelCache {
	cacheMap := make(map[disgord.Snowflake]disgord.Snowflake)
	return &VoiceChannelCache{
		Cache: cacheMap,
	}
}

var vcCache = NewVoiceChannelCache()

// UpdateCache updates the voicechannel cache based upon the set channel id on voice state updates
// A channelID of 0 means a user left, in that case remove them from the cache
func (vcc *VoiceChannelCache) UpdateCache(chID disgord.Snowflake, uID disgord.Snowflake) {
	switch chID {
	case 0:
		delete(vcCache.Cache, uID)
	default:
		vcCache.Cache[uID] = chID
	}
}

// Using this for access to the global clients FOR NOW as passing it through the handlers has proven tricky
// TO-DO find a solution to get rid of the global variables, including the client

// RespondToCommand delegates actions when commands are issued
func RespondToCommand(s disgord.Session, data *disgord.MessageCreate) {
	cec := NewCommandEventClient(data.Message, disgordGlobalClient)
	command, args := cec.DisectCommand()

	fmt.Printf("\nvUserID: %+v\n", data.Message.Author.ID)

	user := data.Message.Author

	fmt.Printf("Command %+v by user %+v | %+v\n", command, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	switch command {
	case "play":
		go processAndPlay(s, data, args[0])
	default:
		cec.Delegate()
	}

}

// RespondToMessage delegates actions when messages are created
func RespondToMessage(s disgord.Session, data *disgord.MessageCreate) {
	user := data.Message.Author

	fmt.Printf("Message %+v by user %+v | %+v\n", data.Message.Content, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	mec := NewMessageEventClient(data.Message, disgordGlobalAPI)
	err := mec.FilterNonModLinks()
	if err != nil {
		fmt.Printf("\nError filtering non-mod link: %+v\n", err)
	}
}

// RespondToReaction delegates actions when reactions are added to messages
func RespondToReaction(s disgord.Session, data *disgord.MessageReactionAdd) {
	userQueryBuilder := disgordGlobalClient.User(data.UserID)

	user, err := userQueryBuilder.Get()

	if err != nil {
		fmt.Printf("\nError getting user: %+v\n", err)
	}
	// fmt.Printf("Message reaction %+v by user %+v | %+v\n", data.PartialEmoji.Name, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	rec := NewReactionEventClient(data.PartialEmoji, data.UserID, data.ChannelID, data.MessageID, disgordGlobalAPI)
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

// RespondToVoiceChannelUpdate updates the server's voice channel cache every time an update is emitted
func RespondToVoiceChannelUpdate(s disgord.Session, data *disgord.VoiceStateUpdate) {
	vcCache.UpdateCache(data.ChannelID, data.UserID)
}

func processAndPlay(s disgord.Session, data *disgord.MessageCreate, arg string) {
	ss := youtube.NewSearchService(ytService)
	ytc := yt.NewYoutubeClient(ss)
	filename, err := ytc.SearchAndDownload(arg)
	if err != nil {
		fmt.Printf("\nERROR WITH FILE: %+v\n", err)
	}

	fmt.Printf("\nFILENAME: %+v\n", filename)

	go deleteMessage(data.Message, 1*time.Second, disgordGlobalClient)
	requestURL := fmt.Sprintf("http://localhost:8080/mp3/%+v", arg)
	encodeSess, err := dca.EncodeFile(requestURL, &dca.EncodeOptions{
		Volume:           256,
		Channels:         2,
		FrameRate:        48000,
		FrameDuration:    20,
		Bitrate:          64,
		Application:      "audio",
		CompressionLevel: 5,
		PacketLoss:       1,
		BufferedFrames:   200, // At 20ms frames that's 2s
		VBR:              true,
		StartTime:        0,
		RawOutput:        true,
		Threads:          0,
	})
	if err != nil {
		fmt.Printf("\nERROR ENCODING: %+v\n", err)
	}

	defer encodeSess.Cleanup()

	vc, err := disgordGlobalClient.VoiceConnectOptions(data.Message.GuildID, data.Message.ChannelID, true, false)
	if err != nil {
		fmt.Printf("\nERROR CONNECTING TO VOICE CHANNEL: %+v\n", err)
		return
	}
	err = vc.StartSpeaking()
	if err != nil {
		fmt.Printf("\nERROR SPEAKING: %+v\n", err)
	}

	go func() {
		err = vc.SendDCA(encodeSess)
		if err != nil {
			fmt.Printf("\nERROR PLAYING DCA: %+v\n", err)
		}
	}()
	// if err := recover(); err != nil {
	// 	log.Println("panic occurred:", err)
	// }

	for {
		time.Sleep(1 * time.Second)
		err = vc.StartSpeaking()
		if err != nil {
			fmt.Printf("\nERROR SPEAKING: %+v\n", err)
		}
		if encodeSess.Running() == false {

			fmt.Printf("\n\nSTOPPING SPEAKING NOW!\n\n")
			vc.StopSpeaking()

			fmt.Printf("\n\nLEAVING VOICE CHAT\n\n")
			vc.Close()

			return
		}
	}
}
