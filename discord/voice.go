package discord

import (
	"fmt"
	"io"
	"time"

	"github.com/andersfylling/disgord"
	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
	"github.com/jonas747/dca"
)

// VoiceEventClient keeps a map record of userID's and the channelID that user currently belongs to
type VoiceEventClient struct {
	Cache         map[disgord.Snowflake]disgord.Snowflake
	disgordClient disgordiface.DisgordClientAPI
}

// NewVoiceChannelCache return a pointer to a new voicechannelcache
func NewVoiceChannelCache(disgordClient disgordiface.DisgordClientAPI) *VoiceEventClient {
	cacheMap := make(map[disgord.Snowflake]disgord.Snowflake)
	return &VoiceEventClient{
		Cache: cacheMap,
	}
}

// UpdateCache updates the voicechannel cache based upon the set channel id on voice state updates
// A channelID of 0 means a user left, in that case remove them from the cache
func (vec *VoiceEventClient) UpdateCache(chID disgord.Snowflake, uID disgord.Snowflake) {
	switch chID {
	case 0:
		delete(vcCache.Cache, uID)
	default:
		vcCache.Cache[uID] = chID
	}
}

// ProcessAndPlay takes a message content string to fetch\encode\play
// audio in the voice channel the author currently resides in
func ProcessAndPlay(gID disgord.Snowflake, uID disgord.Snowflake, arg string, disgordClientAPI disgordiface.DisgordClientAPI) {

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

	vc, err := disgordClientAPI.VoiceConnectOptions(gID, vcCache.Cache[uID], true, false)
	if err != nil {
		fmt.Printf("\nERROR CONNECTING TO VOICE CHANNEL: %+v\n", err)
		// return
	}
	err = vc.StartSpeaking()
	if err != nil {
		fmt.Printf("\nERROR SPEAKING: %+v\n", err)
	}

	ticker := time.NewTicker(20 * time.Millisecond)
	done := make(chan bool)
	eofChannel := make(chan bool)

	go func() {
		defer encodeSess.Cleanup()
		for {
			select {
			case <-done:
				err := vc.StopSpeaking()
				if err != nil && err != io.EOF {
					fmt.Printf("\nERROR STOPPING TALKING: %+v\n", err)
				}
				time.Sleep(1 * time.Second)
				err = vc.Close()
				if err != nil && err != io.EOF {
					fmt.Printf("\nERROR LEAVING VC: %+v\n", err)
				}
				fmt.Println("\nLEAVING VOICE CHANNEL")
				return
			case <-ticker.C:
				nextFrame, err := encodeSess.OpusFrame()
				if err != nil && err != io.EOF {
					fmt.Printf("\nERROR PLAYING DCA: %+v\n", err)
				}
				if err == io.EOF {
					fmt.Println("\nPLAYBACK FINISHED")
					eofChannel <- true
				}
				err = vc.SendOpusFrame(nextFrame)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			select {
			case <-eofChannel:
				done <- true
				return
			}
		}
	}()
}
