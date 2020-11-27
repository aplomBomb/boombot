package discord

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/api/youtube/v3"

	"github.com/andersfylling/disgord"
	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
	"github.com/jonas747/dca"
)

// PlayingDetails contains video data for fetched youtube songs/videos
type PlayingDetails struct {
	Snippet        *youtube.VideoSnippet
	ContentDetails *youtube.VideoContentDetails
	Statistics     *youtube.VideoStatistics
}

// Queue defines the data neccessary for the bot to track users/songs and where to play them
type Queue struct {
	UserQueue               map[disgord.Snowflake][]string
	VoiceCache              map[disgord.Snowflake]disgord.Snowflake
	GuildID                 disgord.Snowflake
	LastMessageUID          disgord.Snowflake
	LastMessageCHID         disgord.Snowflake
	NowPlayingUID           disgord.Snowflake
	LastPlayingUID          disgord.Snowflake
	NowPlayingURL           string
	Next                    chan bool
	Stop                    chan bool
	Shuffle                 chan bool
	ChannelHop              chan disgord.Snowflake
	CurrentlyPlayingDetails PlayingDetails
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
}

// UpdateVoiceCache updates the voicechannel cache based upon the set channel id on voice state updates
// A channelID of 0 means a user left, in that case remove them from the cache
// If the user has a queue list when removed(leaves voice chat entirely) remove their queue entries from the global queue as well
func (q *Queue) UpdateVoiceCache(chID disgord.Snowflake, uID disgord.Snowflake) {
	switch chID {
	case 0:

		delete(q.VoiceCache, uID)
		for k := range q.UserQueue {
			if uID == k {
				delete(q.UserQueue, k)
			}
		}
	default:
		q.VoiceCache[uID] = chID
	}
}

// ListenAndProcessQueue takes a message content string to fetch\encode\play
// audio in the voice channel the author currently resides in
func (q *Queue) ListenAndProcessQueue(disgordClientAPI disgordiface.DisgordClientAPI, youtubeVideosListCall *youtube.VideosListCall) {
	wg := sync.WaitGroup{}
	vc, err := disgordClientAPI.VoiceConnectOptions(q.GuildID, 640284178755092505, true, false)
	if err != nil {
		fmt.Printf("\nERROR: %+v\n", err)
	}
	for {
		time.Sleep(3 * time.Second)
		fmt.Printf("\nNowPlayingUID: %+v | NowPlayingURL: %+v\n", q.NowPlayingUID, q.NowPlayingURL)
		if len(q.UserQueue) > 0 {
			fmt.Println("\nQueues: ", len(q.UserQueue))
			wg.Add(1)
			q.setNowPlaying()
			requestURL := ""
			requestURL = fmt.Sprintf("http://localhost:8080/mp3/%+v", q.UserQueue[q.NowPlayingUID][0])
			fmt.Println("\nURL: ", q.NowPlayingURL)
			fields := strings.Split(q.UserQueue[q.NowPlayingUID][0], "=")
			id := fields[1]
			fmt.Println("\nID: ", id)
			call := youtubeVideosListCall.Id(id)
			resp, err := call.Do()
			if err != nil {
				fmt.Println("\nERROR FETCHING VID DEETZ: ", err)
			}

			q.CurrentlyPlayingDetails.Snippet = resp.Items[0].Snippet
			q.CurrentlyPlayingDetails.ContentDetails = resp.Items[0].ContentDetails
			q.CurrentlyPlayingDetails.Statistics = resp.Items[0].Statistics

			if len(resp.Items) == 0 {
				fmt.Println("\nNo data retrieved from Youtube|Skipping...")
				q.RemoveQueueEntry()
				wg.Done()
				continue
			}
			es, err := q.GetEncodeSession(requestURL)
			if err != nil {
				fmt.Printf("\nERROR ENCODING: %+v\n", err)
			}
			esData, err := es.ReadFrame()
			if err != nil {
				fmt.Println("\nError: ", err)
			}
			if esData == nil {
				fmt.Println("\nNo audio data|Skipping...")
				q.RemoveQueueEntry()
				wg.Done()
				continue
			}
			vc, err = q.establishVoiceConnection(vc, disgordClientAPI, q.VoiceCache[739154323015204935], q.VoiceCache[q.NowPlayingUID])
			if err != nil {
				fmt.Printf("\nERROR: %+v\n", err)
			}

			fmt.Printf("\nNowPlayingUID: %+v | NowPlayingURL: %+v\n", q.NowPlayingUID, q.NowPlayingURL)

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
		if len(q.UserQueue) == 0 {
			q.NowPlayingUID = 0
			q.CurrentlyPlayingDetails = PlayingDetails{}
			q.NowPlayingURL = ""
		}
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
	referenceEntry := PlayingDetails{}
	for {
		time.Sleep(2 * time.Second)
		// fmt.Println("\nNOW PLAYING: ", q.UserQueue[q.NowPlayingUID])
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
			if referenceEntry != q.CurrentlyPlayingDetails {
				// nextRequesteeName := "**Open Queue**"
				requesteeName, err := disgordClient.User(q.NowPlayingUID).Get()
				if err != nil {
					fmt.Println("\n", err)
				}
				avatarURL, err := requesteeName.AvatarURL(64, true)
				if err != nil {
					fmt.Println("\n", err)
				}
				likeStr := strconv.FormatUint(q.CurrentlyPlayingDetails.Statistics.LikeCount, 10)
				dislikeStr := strconv.FormatUint(q.CurrentlyPlayingDetails.Statistics.DislikeCount, 10)
				timeFields := strings.Split(q.CurrentlyPlayingDetails.ContentDetails.Duration, "PT")
				disgordClient.SendMsg(
					779836590503624734,
					&disgord.CreateMessageParams{
						Embed: &disgord.Embed{
							Title:       q.CurrentlyPlayingDetails.Snippet.Title,
							URL:         q.UserQueue[q.NowPlayingUID][0],
							Description: q.CurrentlyPlayingDetails.Snippet.Description,
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
								&disgord.EmbedField{
									Name:   "Duration",
									Value:  timeFields[1],
									Inline: true,
								},
								&disgord.EmbedField{
									Name:   "Upvotes",
									Value:  likeStr,
									Inline: true,
								},
								&disgord.EmbedField{
									Name:   "Downvotes",
									Value:  dislikeStr,
									Inline: true,
								},
							},
							Image: &disgord.EmbedImage{
								URL:    q.CurrentlyPlayingDetails.Snippet.Thumbnails.High.Url,
								Height: 128,
								Width:  128,
							},
							Footer: &disgord.EmbedFooter{
								IconURL: avatarURL,
								Text:    fmt.Sprintf("%+v's queue: %+v", requesteeName.Username, strconv.Itoa(len(q.UserQueue[q.NowPlayingUID]))),
							},
						},
					},
				)
				if len(q.UserQueue) > 0 {
					referenceEntry = q.CurrentlyPlayingDetails
				}
			}
		}
		if len(msgs) > 1 {

			go deleteMessage(msgs[1], 1*time.Second, disgordClient)

		}
		if q.NowPlayingUID == 0 {
			for k := range msgs {
				go deleteMessage(msgs[k], 1*time.Second, disgordClient)
			}
		}
	}
}

// RemoveQueueEntry removes the last queue entry and deletes the map if string slice is empty
// This is insane looking/literally makes my eyes glaze over looking at it
// Should devise a more user(reader)-friendly solution
func (q *Queue) RemoveQueueEntry() {

	nowPlayingID := disgord.Snowflake(0)

	for k := range q.UserQueue {
		if k == q.NowPlayingUID {
			nowPlayingID = q.NowPlayingUID
		}
	}

	if nowPlayingID != 0 {
		copy(q.UserQueue[q.NowPlayingUID][0:], q.UserQueue[q.NowPlayingUID][0+1:])
		q.UserQueue[q.NowPlayingUID][len(q.UserQueue[q.NowPlayingUID])-1] = ""
		q.UserQueue[q.NowPlayingUID] = q.UserQueue[q.NowPlayingUID][:len(q.UserQueue[q.NowPlayingUID])-1]
		if len(q.UserQueue[q.NowPlayingUID]) <= 0 {
			delete(q.UserQueue, q.NowPlayingUID)
		}
	}

	q.LastPlayingUID = q.NowPlayingUID
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
	q.LastPlayingUID = q.NowPlayingUID
}

func (q *Queue) stopPlaybackAndTalking(vc disgord.VoiceConnection, es *dca.EncodeSession) {
	err := es.Stop()
	if err != nil {
		fmt.Printf("\nERROR STOPPING ENCODING: %+v", err)
	}
	err = vc.StopSpeaking()
	if err != nil && err != io.EOF {
		fmt.Printf("\nERROR STOPPING TALKING: %+v\n", err)
	}
}

func (q *Queue) establishVoiceConnection(prevVC disgord.VoiceConnection, client disgordiface.DisgordClientAPI, botChannelID disgord.Snowflake, requesteeChannelID disgord.Snowflake) (disgord.VoiceConnection, error) {
	if botChannelID == 0 {
		vc, err := client.VoiceConnectOptions(q.GuildID, requesteeChannelID, true, false)
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
		prevVC.Close()
		newVC, err := client.VoiceConnectOptions(q.GuildID, requesteeChannelID, true, false)
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
		err := prevVC.StartSpeaking()
		if err != nil {
			return nil, err
		}

		return prevVC, nil
	}
	return prevVC, nil
}

func (q *Queue) setNowPlaying() {
	lastUID := q.LastPlayingUID
	nextUID := disgord.Snowflake(0)

	fmt.Println("\nLastPlayingUID: ", lastUID)

	// When there is more than one queue, we dont want to play the same user's queue twice in a row
	// Collect all the queue id's that aren't the last one
	uidbucket := []disgord.Snowflake{}
	if len(q.UserQueue) > 1 {
		for k := range q.UserQueue {
			if k != lastUID {
				uidbucket = append(uidbucket, k)
			}
		}
	} else {
		for k := range q.UserQueue {
			lastUID = k
		}
		uidbucket = append(uidbucket, lastUID)
	}

	if len(uidbucket) > 1 {
		// And if there's more than one, pick one at random
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		nextUID = uidbucket[r.Intn(len(uidbucket))]
	} else {
		nextUID = uidbucket[0]
	}
	q.NowPlayingUID = nextUID
	q.NowPlayingURL = q.UserQueue[q.NowPlayingUID][0]
}

// RemoveQueueByID removes a user's queue via their userID
func (q *Queue) RemoveQueueByID(id disgord.Snowflake) {
	delete(q.UserQueue, id)
}
