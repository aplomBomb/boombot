package discord

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
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
	NowPlayingUID   disgord.Snowflake
	LastPlayingUID  disgord.Snowflake
	NextPlayingUID  disgord.Snowflake
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
		LastPlayingUID:  disgord.Snowflake(0),
		NowPlayingUID:   disgord.Snowflake(0),
		NextPlayingUID:  disgord.Snowflake(0),
		NowPlayingURL:   "",
		Next:            make(chan bool, 1),
		Stop:            make(chan bool, 1),
		Shuffle:         make(chan bool, 1),
		ChannelHop:      make(chan disgord.Snowflake, 1),
	}
}

// UpdateUserQueueState updates the UserQueue cache on single song play requests
func (q *Queue) UpdateUserQueueState(chID disgord.Snowflake, uID disgord.Snowflake, arg string) {
	q.LastMessageUID = uID
	q.LastMessageCHID = chID
	if q.UserQueue[uID] == nil {
		q.UserQueue[uID] = []string{arg}
	} else {
		q.UserQueue[uID] = append(q.UserQueue[uID], arg)
	}
	// ================ I dont know if i should be doing this here =========================
	// if this works, turn it into a method
	uidbucket := []disgord.Snowflake{}
	i := 0
	for k := range q.UserQueue {
		uidbucket = append(uidbucket, k)
	}
	for k, v := range uidbucket {
		if v == q.LastPlayingUID {
			i = k - 1
		}
	}
	if i < 0 {
		i = len(uidbucket) - 1
	}
	q.NextPlayingUID = uidbucket[i]

	// =====================================================================================
}

// UpdateUserQueueStateBulk updates the UserQueue and GlobalQueue cache for playlist requests
// AKA 'Playloads' lolol
func (q *Queue) UpdateUserQueueStateBulk(chID disgord.Snowflake, uID disgord.Snowflake, args []string) {
	q.LastMessageUID = uID
	q.LastMessageCHID = chID
	if q.UserQueue[uID] == nil {
		q.UserQueue[uID] = args
	} else {
		q.UserQueue[uID] = append(q.UserQueue[uID], args...)
	}
	// ================ I dont know if i should be doing this here =========================
	// if this works, turn it into a method
	uidbucket := []disgord.Snowflake{}
	i := 0
	for k := range q.UserQueue {
		uidbucket = append(uidbucket, k)
	}
	for k, v := range uidbucket {
		if v == q.LastPlayingUID {
			i = k - 1
		}
	}
	if i < 0 {
		i = len(uidbucket) - 1
	}
	q.NextPlayingUID = uidbucket[i]

	// =====================================================================================
}

// UpdateVoiceCache updates the voicechannel cache based upon the set channel id on voice state updates
// A channelID of 0 means a user left, in that case remove them from the cache
// If the user has a queue list when removed(leaves voice chat entirely) remove their queue entries from the global queue as well
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
	vc, err := disgordClientAPI.VoiceConnectOptions(q.GuildID, 640284178755092505, true, false)
	if err != nil {
		fmt.Printf("\nERROR: %+v\n", err)
	}
	for {
		time.Sleep(3 * time.Second)
		if len(q.UserQueue) > 0 {
			wg.Add(1)
			requestURL := ""
			requestURL = fmt.Sprintf("http://localhost:8080/mp3/%+v", q.UserQueue[q.NextPlayingUID][0])
			fmt.Printf("\nPLAYING FROM USER PLAYLIST FROM LAST COMMAND\n")
			q.NowPlayingURL = q.UserQueue[q.NextPlayingUID][0]
			q.NowPlayingSync()
			fmt.Printf("\nUSER CURRENT: %+v\n", q.NowPlayingUID)
			es, err := q.GetEncodeSession(requestURL)
			if err != nil {
				fmt.Printf("\nERROR ENCODING: %+v\n", err)
			}

			vc, err = q.establishVoiceConnection(vc, disgordClientAPI, q.VoiceCache[739154323015204935], q.VoiceCache[q.NowPlayingUID])
			if err != nil {
				fmt.Printf("\nERROR: %+v\n", err)
			}

			// Ticker needed for smooth opus frame delivery to prevent playback stuttering
			ticker := time.NewTicker(20 * time.Millisecond)
			done := make(chan bool)
			eofChannel := make(chan bool)
			// Goroutine just cycles through the opusFrames produced from the encoding process
			// The channels allow for realtime interaction/playback control from events triggered by users
			go func(waitGroup *sync.WaitGroup) {
				defer es.Cleanup()
				defer ticker.Stop()
				defer waitGroup.Done()
				for {
					select {
					case <-q.Shuffle:
						q.stopPlaybackAndTalking(vc, es)
						q.ShuffleQueue()
						time.Sleep(1 * time.Second)
						return
					case <-q.Stop:
						q.stopPlaybackAndTalking(vc, es)
						q.EmptyQueue()
						time.Sleep(1 * time.Second)
						return
					case <-q.Next:
						q.stopPlaybackAndTalking(vc, es)
						q.RemoveQueueEntry()
						time.Sleep(1 * time.Second)
						return
					case <-done:
						q.stopPlaybackAndTalking(vc, es)
						q.RemoveQueueEntry()
						time.Sleep(1 * time.Second)

						return
					case channelID := <-q.ChannelHop:
						vc.StopSpeaking()
						vc, err = q.establishVoiceConnection(vc, disgordClientAPI, 0, channelID)
						if err != nil {
							fmt.Printf("\nERROR: %+v\n", err)
						}
						vc.StartSpeaking()
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
		// TO-DO===============================================================
		// I need to respond to a message create event rather than pinging discord's api every two seconds
		//(2 seconds is the tightest interval before being rate-limited)
		msgs, err := disgordClient.Channel(779836590503624734).GetMessages(&disgord.GetMessagesParams{
			Limit: 10,
		})
		if err != nil {
			fmt.Printf("\nCOULD NOT GET MESSAGES FROM JUKEBOX CHANNEL: %+v", err)
		}
		if len(q.UserQueue) > 0 && q.NowPlayingUID != 0 {
			if referenceEntry != q.UserQueue[q.NowPlayingUID][0] {

				requesteeName, err := disgordClient.User(q.NowPlayingUID).Get()
				if err != nil {
					fmt.Println("\n", err)
				}
				avatarURL, err := requesteeName.AvatarURL(64, true)
				if err != nil {
					fmt.Println("\n", err)
				}
				disgordClient.SendMsg(
					779836590503624734,
					&disgord.CreateMessageParams{
						Content: q.UserQueue[q.NowPlayingUID][0],
					},
				)
				disgordClient.SendMsg(
					779836590503624734,
					&disgord.CreateMessageParams{
						Embed: &disgord.Embed{
							Thumbnail: &disgord.EmbedThumbnail{
								URL:    avatarURL,
								Height: 64,
								Width:  64,
							},
							Fields: []*disgord.EmbedField{
								&disgord.EmbedField{
									Name:  "Requested by",
									Value: requesteeName.Username,
								},
							},
							Footer: &disgord.EmbedFooter{
								Text: fmt.Sprintf("%+v's queue: %+v", requesteeName.Username, strconv.Itoa(len(q.UserQueue[q.NowPlayingUID]))),
							},
						},
					},
				)
				if len(q.UserQueue) > 0 {
					referenceEntry = q.UserQueue[q.NowPlayingUID][0]
				}
			}
		}
		if len(msgs) > 2 {
			for k := range msgs {
				if k > 1 {
					go deleteMessage(msgs[k], 1*time.Second, disgordClient)
				}
			}
		}
	}
}

// RemoveQueueEntry removes the last queue entry and deletes the map if string slice is empty
// This is insane looking/literally makes my eyes glaze over looking at it
// Should devise a more user(reader)-friendly solution
func (q *Queue) RemoveQueueEntry() {
	copy(q.UserQueue[q.NowPlayingUID][0:], q.UserQueue[q.NowPlayingUID][0+1:])
	q.UserQueue[q.NowPlayingUID][len(q.UserQueue[q.NowPlayingUID])-1] = ""
	q.UserQueue[q.NowPlayingUID] = q.UserQueue[q.NowPlayingUID][:len(q.UserQueue[q.NowPlayingUID])-1]
	if len(q.UserQueue[q.NowPlayingUID]) <= 0 {
		delete(q.UserQueue, q.NowPlayingUID)
		q.NowPlayingUID = 0
	}
}

// ShuffleQueue reorganizes the order of the queue entries for randomized playback
func (q *Queue) ShuffleQueue() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(q.UserQueue[q.NowPlayingUID]), func(i, j int) {
		q.UserQueue[q.NowPlayingUID][i], q.UserQueue[q.NowPlayingUID][j] = q.UserQueue[q.NowPlayingUID][j], q.UserQueue[q.NowPlayingUID][i]
	})
}

// EmptyQueue deletes user's queue map
func (q *Queue) EmptyQueue() {
	delete(q.UserQueue, q.NowPlayingUID)
	q.NowPlayingUID = 0
}

// NowPlayingSync keeps the NowPlayingUID updated with the ID of the user who's song is currently playing
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
	q.NowPlayingUID = currentUID
}

func (q *Queue) stopPlaybackAndTalking(vc disgord.VoiceConnection, es *dca.EncodeSession) {
	err := es.Stop()
	if err != nil {
		fmt.Printf("\nERROR STOPPING ENCODING: %+v", err)
	}
	fmt.Println("\nSKIPPING QUEUE ENTRY")
	err = vc.StopSpeaking()
	if err != nil && err != io.EOF {
		fmt.Printf("\nERROR STOPPING TALKING: %+v\n", err)
	}
}

func (q *Queue) establishVoiceConnection(prevVC disgord.VoiceConnection, client disgordiface.DisgordClientAPI, botChannelID disgord.Snowflake, requesteeChannelID disgord.Snowflake) (disgord.VoiceConnection, error) {
	if botChannelID == 0 {
		fmt.Println("\nJoining requestee's voice channel")
		vc, err := client.VoiceConnectOptions(q.GuildID, requesteeChannelID, true, false)
		// queueCounter := 0
		if err != nil {
			return nil, err
		}

		err = vc.StartSpeaking()
		if err != nil {
			return nil, err
		}

		prevVC = vc
		return vc, nil
	}
	if botChannelID != requesteeChannelID {
		fmt.Println("\nJoining requestee's voice channel")
		prevVC.Close()
		newVC, err := client.VoiceConnectOptions(q.GuildID, requesteeChannelID, true, false)
		// queueCounter := 0
		if err != nil {
			return nil, err
		}

		err = newVC.StartSpeaking()
		if err != nil {
			return nil, err
		}

		return newVC, nil
	}
	if botChannelID == requesteeChannelID {
		fmt.Println("\nPlaying next song from playlist in same channel")

		err := prevVC.StartSpeaking()
		if err != nil {
			return nil, err
		}

		return prevVC, nil
	}

	return prevVC, nil
}
