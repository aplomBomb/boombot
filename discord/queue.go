package discord

import (
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/andersfylling/disgord"
	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
	"github.com/jonas747/dca"
)

// Queue defines the data neccessary for the bot to track users/songs and where to play them
type Queue struct {
	UserQueue       []string
	VoiceCache      map[disgord.Snowflake]disgord.Snowflake
	GuildID         disgord.Snowflake
	LastMessageUID  disgord.Snowflake
	LastMessageCHID disgord.Snowflake
	Next            chan bool
	Stop            chan bool
	Shuffle         chan bool
}

// NewQueue returns a new Queue instance
func NewQueue(gID disgord.Snowflake) *Queue {
	return &Queue{
		UserQueue:       []string{},
		VoiceCache:      map[disgord.Snowflake]disgord.Snowflake{},
		GuildID:         gID,
		LastMessageUID:  disgord.Snowflake(0),
		LastMessageCHID: disgord.Snowflake(0),
		Next:            make(chan bool),
		Stop:            make(chan bool),
		Shuffle:         make(chan bool),
	}
}

// UpdateQueueState updates the Queue cache on single song play requests
func (q *Queue) UpdateQueueState(chID disgord.Snowflake, uID disgord.Snowflake, arg string) {
	q.LastMessageUID = uID
	q.LastMessageCHID = chID
	q.UserQueue = append(q.UserQueue, arg)
}

// UpdateQueueStateBulk updates the Queue cache for playlist payload requests
func (q *Queue) UpdateQueueStateBulk(chID disgord.Snowflake, uID disgord.Snowflake, args []string) {
	q.LastMessageUID = uID
	q.LastMessageCHID = chID
	q.UserQueue = append(q.UserQueue, args...)
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
	wg := sync.WaitGroup{}
	for {
		time.Sleep(3 * time.Second)
		if len(q.UserQueue) > 0 {
			wg.Add(1)
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
			vc, err := disgordClientAPI.VoiceConnectOptions(q.GuildID, q.VoiceCache[q.LastMessageUID], true, false)
			if err != nil {
				fmt.Printf("\nERROR CONNECTING TO VOICE CHANNEL: %+v\n", err)
				msg, err := disgordClientAPI.SendMsg(q.LastMessageCHID, disgord.Message{
					Content: "**You need to be in a voice channel for me to play!**",
				})
				if err != nil {
					fmt.Printf("\nERROR SENDING NO CHANNELID MESSSAGE: %+v\n", err)
				}
				//Delete url from queue slice
				copy(q.UserQueue[0:], q.UserQueue[0+1:])
				q.UserQueue[len(q.UserQueue)-1] = ""
				q.UserQueue = q.UserQueue[:len(q.UserQueue)-1]
				go deleteMessage(msg, 10*time.Second, disgordClientAPI)
				wg.Done()
				continue
			}
			err = vc.StartSpeaking()
			if err != nil {
				fmt.Printf("\nERROR SPEAKING: %+v\n", err)
			}
			// Ticker needed for smooth opus frame delivery to prevent playback stuttering
			ticker := time.NewTicker(20 * time.Millisecond)
			done := make(chan bool)
			eofChannel := make(chan bool)
			go func(waitGroup *sync.WaitGroup) {
				defer encodeSess.Cleanup()
				defer ticker.Stop()
				defer waitGroup.Done()
				for {
					select {
					case <-q.Shuffle:
						err := encodeSess.Stop()
						if err != nil {
							fmt.Printf("\nERROR STOPPING ENCODING: %+v", err)
						}
						err = vc.StopSpeaking()
						if err != nil && err != io.EOF {
							fmt.Printf("\nERROR STOPPING TALKING: %+v\n", err)
						}
						time.Sleep(1 * time.Second)
						err = vc.Close()
						if err != nil && err != io.EOF {
							fmt.Printf("\nERROR LEAVING VC: %+v\n", err)
						}
						fmt.Println("\nLEAVING VOICE CHANNEL")
						q.ShuffleQueue()
						return
					case <-q.Stop:
						err := encodeSess.Stop()
						if err != nil {
							fmt.Printf("\nERROR STOPPING ENCODING: %+v", err)
						}
						err = vc.StopSpeaking()
						if err != nil && err != io.EOF {
							fmt.Printf("\nERROR STOPPING TALKING: %+v\n", err)
						}
						time.Sleep(1 * time.Second)
						err = vc.Close()
						if err != nil && err != io.EOF {
							fmt.Printf("\nERROR LEAVING VC: %+v\n", err)
						}
						fmt.Println("\nLEAVING VOICE CHANNEL")
						q.EmptyQueue()
						return
					case <-q.Next:
						err := encodeSess.Stop()
						if err != nil {
							fmt.Printf("\nERROR STOPPING ENCODING: %+v", err)
						}
						fmt.Println("SKIPPING QUEUE ENTRY")
						err = vc.StopSpeaking()
						if err != nil && err != io.EOF {
							fmt.Printf("\nERROR STOPPING TALKING: %+v\n", err)
						}
						time.Sleep(1 * time.Second)
						err = vc.Close()
						if err != nil && err != io.EOF {
							fmt.Printf("\nERROR LEAVING VC: %+v\n", err)
						}
						fmt.Println("\nLEAVING VOICE CHANNEL")
						q.RemoveLastQueueEntry()
						return
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
						q.RemoveLastQueueEntry()
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
						vc.SendOpusFrame(nextFrame)
					}
				}
			}(&wg)
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
		wg.Wait()
	}
}

// ManageJukebox scans the userqueue and updates the embed in the jukebox designated channel
func (q *Queue) ManageJukebox(disgordClient disgordiface.DisgordClientAPI) {
	referenceEntry := ""
	for {
		time.Sleep(5 * time.Second)
		msgs, err := disgordClient.Channel(779836590503624734).GetMessages(&disgord.GetMessagesParams{
			Limit: 10,
		})
		if err != nil {
			fmt.Printf("\nCOULD NOT GET MESSAGES FROM JUKEBOX CHANNEL: %+v", err)
		}
		if len(q.UserQueue) > 0 {
			if referenceEntry != q.UserQueue[0] {
				disgordClient.SendMsg(
					779836590503624734,
					&disgord.CreateMessageParams{
						Content: q.UserQueue[0],
					},
				)
				if len(q.UserQueue) > 0 {
					referenceEntry = q.UserQueue[0]
				}
			}
		}
		if len(msgs) > 1 {
			go deleteMessage(msgs[1], 1*time.Millisecond, disgordClient)
		}
	}
}

// RemoveLastQueueEntry does what it says, removes the last entry in the queue
func (q *Queue) RemoveLastQueueEntry() {
	copy(q.UserQueue[0:], q.UserQueue[0+1:])
	q.UserQueue[len(q.UserQueue)-1] = ""
	q.UserQueue = q.UserQueue[:len(q.UserQueue)-1]
}

// ShuffleQueue randomizes the order of the queue entries for unpredictable playback
func (q *Queue) ShuffleQueue() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(q.UserQueue), func(i, j int) { q.UserQueue[i], q.UserQueue[j] = q.UserQueue[j], q.UserQueue[i] })
}

// EmptyQueue reverts queue to it's zero value state
func (q *Queue) EmptyQueue() {
	q.UserQueue = []string{}
}
