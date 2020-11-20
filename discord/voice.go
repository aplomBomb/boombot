package discord

import (
	"fmt"
	"io"
	"time"

	"github.com/andersfylling/disgord"
	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
	"github.com/jonas747/dca"
)

// Queue defines the data neccessary for the bot to track users/songs and where to play them
type Queue struct {
	UserQueue  map[disgord.Snowflake][]string
	VoiceCache map[disgord.Snowflake]disgord.Snowflake
	GuildID    disgord.Snowflake
}

// NewQueue returns a new Queue instance
func NewQueue(gID disgord.Snowflake) *Queue {
	return &Queue{
		UserQueue:  map[disgord.Snowflake][]string{},
		VoiceCache: map[disgord.Snowflake]disgord.Snowflake{},
		GuildID:    gID,
	}
}

// UpdateVoiceCache updates the voicechannel cache based upon the set channel id on voice state updates
// A channelID of 0 means a user left, in that case remove them from the cache
func (q *Queue) UpdateVoiceCache(chID disgord.Snowflake, uID disgord.Snowflake) {
	switch chID {
	case 0:
		delete(q.VoiceCache, uID)
	default:
		q.VoiceCache[uID] = chID
	}
}

// ListenAndProcessQueue takes a message content string to fetch\encode\play
// audio in the voice channel the author currently resides in
func (q *Queue) ListenAndProcessQueue(disgordClientAPI disgordiface.DisgordClientAPI) {
	requestURL := fmt.Sprintf("http://localhost:8080/mp3/%+v", q.UserQueue[0])
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
	vc, err := disgordClientAPI.VoiceConnectOptions(q.GuildID, q.VoiceCache[uID], true, false)
	if err != nil {
		fmt.Printf("\nERROR CONNECTING TO VOICE CHANNEL: %+v\n", err)
		msg, err := disgordClientAPI.SendMsg(chID, disgord.Message{
			Content: "**You need to be in a voice channel for me to play!**",
		})
		if err != nil {
			fmt.Printf("\nERROR SENDING NO CHANNELID MESSSAGE: %+v\n", err)
		}
		go deleteMessage(msg, 10*time.Second, disgordClientAPI)
		return
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
		defer ticker.Stop()
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
