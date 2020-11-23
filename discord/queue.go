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
	UserQueue       map[disgord.Snowflake][]string
	VoiceCache      map[disgord.Snowflake]disgord.Snowflake
	GuildID         disgord.Snowflake
	LastMessageUID  disgord.Snowflake
	LastMessageCHID disgord.Snowflake
	NowPlayinguID   disgord.Snowflake
	NowPlayingURL   string
	Next            chan bool
	Stop            chan bool
	Shuffle         chan bool
	ChannelHop      chan disgord.Snowflake
}

// NewQueue returns a new Queue instance
func NewQueue(gID disgord.Snowflake) *Queue {
	return &Queue{
		UserQueue:       map[disgord.Snowflake][]string{},
		VoiceCache:      map[disgord.Snowflake]disgord.Snowflake{},
		GuildID:         gID,
		LastMessageUID:  disgord.Snowflake(0),
		LastMessageCHID: disgord.Snowflake(0),
		NowPlayinguID:   disgord.Snowflake(0),
		NowPlayingURL:   "",
		Next:            make(chan bool, 1),
		Stop:            make(chan bool, 1),
		Shuffle:         make(chan bool, 1),
		ChannelHop:      make(chan disgord.Snowflake, 1),
	}
}

// UpdateQueueState updates the Queue cache on single song play requests
func (q *Queue) UpdateQueueState(chID disgord.Snowflake, uID disgord.Snowflake, arg string) {
	q.LastMessageUID = uID
	q.LastMessageCHID = chID
	if q.UserQueue[uID] == nil {
		q.UserQueue[uID] = []string{arg}
	} else {
		q.UserQueue[uID] = append(q.UserQueue[uID], arg)
	}
}

// UpdateQueueStateBulk updates the Queue cache for playlist requests
// AKA 'Playloads' lolol
func (q *Queue) UpdateQueueStateBulk(chID disgord.Snowflake, uID disgord.Snowflake, args []string) {
	q.LastMessageUID = uID
	q.LastMessageCHID = chID
	if q.UserQueue[uID] == nil {
		q.UserQueue[uID] = args
	} else {
		q.UserQueue[uID] = append(q.UserQueue[uID], args...)
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
	wg := sync.WaitGroup{}
	for {
		time.Sleep(3 * time.Second)
		// fmt.Printf("\nCURRENTUSERQUEUE: %+v\n", q.UserQueue)
		// fmt.Printf("\nLASTUSERCOMMANDQUEUE: %+v\n", q.UserQueue[q.LastMessageUID])
		if len(q.UserQueue) > 0 {
			wg.Add(1)
			requestURL := ""
			// songURL := ""
			if q.NowPlayinguID == 0 {
				requestURL = fmt.Sprintf("http://localhost:8080/mp3/%+v", q.UserQueue[q.LastMessageUID][0])
				fmt.Printf("\nPLAYING FROM USER PLAYLIST FROM LAST COMMAND\n")
				q.NowPlayingURL = q.UserQueue[q.LastMessageUID][0]
			} else {
				requestURL = fmt.Sprintf("http://localhost:8080/mp3/%+v", q.UserQueue[q.NowPlayinguID][0])
				fmt.Printf("\nCONTINUING FROM USER PLAYLIST\n")
				q.NowPlayingURL = q.UserQueue[q.NowPlayinguID][0]
			}
			q.NowPlayingSync()
			fmt.Printf("\nUSER CURRENT: %+v\n", q.NowPlayinguID)
			es, err := q.GetEncodeSession(requestURL)
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
				q.RemoveQueueEntry()
				go deleteMessage(msg, 5*time.Second, disgordClientAPI)
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
				defer es.Cleanup()
				defer ticker.Stop()
				defer waitGroup.Done()
				for {
					select {
					case <-q.Shuffle:
						err := es.Stop()
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
						err := es.Stop()
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
						err := es.Stop()
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
						q.RemoveQueueEntry()
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
						q.RemoveQueueEntry()
						return
					case channelID := <-q.ChannelHop:
						fmt.Printf("\nSong requester jumped to %+v!", channelID)
					case <-ticker.C:
						nextFrame, err := es.OpusFrame()
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

// GetEncodeSession returns a dca encoded session
func (q *Queue) GetEncodeSession(url string) (*dca.EncodeSession, error) {
	encodeSess, err := dca.EncodeFile(url, &dca.EncodeOptions{
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
		return nil, err
	}
	return encodeSess, nil
}

// ManageJukebox scans the userqueue and updates the embed in the jukebox designated channel
func (q *Queue) ManageJukebox(disgordClient disgordiface.DisgordClientAPI) {
	referenceEntry := ""
	for {
		time.Sleep(2 * time.Second)
		// fmt.Printf("\nCURRENTUSERQUEUE: %+v\n", q.UserQueue)
		// fmt.Printf("\nLASTUSERCOMMANDQUEUE: %+v\n", q.UserQueue[q.LastMessageUID])
		msgs, err := disgordClient.Channel(779836590503624734).GetMessages(&disgord.GetMessagesParams{
			Limit: 10,
		})
		if err != nil {
			fmt.Printf("\nCOULD NOT GET MESSAGES FROM JUKEBOX CHANNEL: %+v", err)
		}
		if len(q.UserQueue) > 0 && q.NowPlayinguID != 0 {
			if referenceEntry != q.UserQueue[q.NowPlayinguID][0] {
				disgordClient.SendMsg(
					779836590503624734,
					&disgord.CreateMessageParams{
						Content: q.UserQueue[q.NowPlayinguID][0],
					},
				)
				if len(q.UserQueue) > 0 {
					referenceEntry = q.UserQueue[q.NowPlayinguID][0]
				}
			}
		}
		if len(msgs) > 1 {
			go deleteMessage(msgs[1], 1*time.Second, disgordClient)
		}
	}
}

// RemoveQueueEntry removes the last queue entry and deletes the map if string slice is empty
func (q *Queue) RemoveQueueEntry() {
	copy(q.UserQueue[q.NowPlayinguID][0:], q.UserQueue[q.NowPlayinguID][0+1:])
	q.UserQueue[q.NowPlayinguID][len(q.UserQueue[q.NowPlayinguID])-1] = ""
	q.UserQueue[q.NowPlayinguID] = q.UserQueue[q.NowPlayinguID][:len(q.UserQueue[q.NowPlayinguID])-1]
	if len(q.UserQueue[q.NowPlayinguID]) <= 0 {
		delete(q.UserQueue, q.NowPlayinguID)
		q.NowPlayinguID = 0
	}
}

// ShuffleQueue randomizes the order of the queue entries for randomized playback
func (q *Queue) ShuffleQueue() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(q.UserQueue[q.NowPlayinguID]), func(i, j int) {
		q.UserQueue[q.NowPlayinguID][i], q.UserQueue[q.NowPlayinguID][j] = q.UserQueue[q.NowPlayinguID][j], q.UserQueue[q.NowPlayinguID][i]
	})
}

// EmptyQueue deletes user's queue map
func (q *Queue) EmptyQueue() {
	delete(q.UserQueue, q.NowPlayinguID)
	q.NowPlayinguID = 0
}

// NowPlayingSync keeps the NowPlayinguID updated with the currently playing queue item
func (q *Queue) NowPlayingSync() {
	currentUID := disgord.Snowflake(0)
	i := 0
	for idKey, stringArr := range q.UserQueue {
		for _, url := range stringArr {
			if url == q.NowPlayingURL {
				currentUID = idKey
			}
		}
		i++
	}
	q.NowPlayinguID = currentUID
}
