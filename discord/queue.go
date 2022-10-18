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
	Pause                   chan bool
	Play                    chan bool
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
		Pause:           make(chan bool, 1),
		Play:            make(chan bool, 1),
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
		q.ShuffleQueueByID(uID)
	} else {
		q.UserQueue[uID] = append(q.UserQueue[uID], args...)
		q.ShuffleNowPlayingQueue()
		q.ShuffleQueueByID(uID)
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

// TriggerNext sends a true boolean value to the queue next channel, skipping whatever queue entry is currently playing
func (q *Queue) TriggerNext() {
	q.Next <- true
}

// TriggerShuffle sends a true boolean value to the queue shuffle channel, shuffling whatever user queue is currently active
func (q *Queue) TriggerShuffle() {
	q.Shuffle <- true
}

// TriggerStop sends a true boolean value to the queue stop channel, stopping whatever is currently playing
func (q *Queue) TriggerStop() {
	q.Stop <- true
}

// TriggerChannelHop sends a channelID to the queue channelhop channel
// For cases when the bot needs to follow a user who's song is currently playing
func (q *Queue) TriggerChannelHop(id disgord.Snowflake) {
	q.ChannelHop <- id
}

// ReturnVoiceCacheEntry returns a voice queue voicecache channel id via a user's id
func (q *Queue) ReturnVoiceCacheEntry(id disgord.Snowflake) disgord.Snowflake {
	return q.VoiceCache[id]
}

// ReturnUserQueue returns the current userqueue from the global queue cache
func (q *Queue) ReturnUserQueue() map[disgord.Snowflake][]string {
	return q.UserQueue
}

// ReturnNowPlayingID returns the nowplayingid from the queue
func (q *Queue) ReturnNowPlayingID() disgord.Snowflake {
	return q.NowPlayingUID
}

// ListenAndProcessQueue takes a message content string to fetch\encode\play
// audio in the voice channel the author currently resides in
func (q *Queue) ListenAndProcessQueue(disgordClientAPI disgordiface.DisgordClientAPI, guild disgordiface.GuildQueryBuilderAPI, ytvlc *youtube.VideosListCall) {
	wg := sync.WaitGroup{}
	vcBuilder := guild.VoiceChannel(915762663752077342)

	vc, err := vcBuilder.Connect(true, false)
	if err != nil {
		fmt.Printf("\nERROR: %+v\n", err)
	}
	for {
		time.Sleep(3 * time.Second)
		if len(q.UserQueue) > 0 {
			fmt.Println("\nQueues: ", len(q.UserQueue))
			wg.Add(1)
			fmt.Printf("\nUpcoming Song/URL: %+v", q.UserQueue[q.NowPlayingUID][0])
			q.setNowPlaying()
			requestURL := ""

			// Localhost address for local testing/development
			// use this address when running the containers independently on the same machine
			// requestURL = fmt.Sprintf("http://localhost:8080/mp3/%+v", q.UserQueue[q.NowPlayingUID][0])

			// yt-api is the name of the intermediary container that fetches youtube audio data for encoding
			// use this address when running the containers together via docker-compose
			requestURL = fmt.Sprintf("http://yt-api:8080/mp3/%+v", q.UserQueue[q.NowPlayingUID][0])

			fields := strings.Split(q.UserQueue[q.NowPlayingUID][0], "=")
			id := fields[1]
			fmt.Printf("SONG_ID: %+v", id)

			call := ytvlc.Id(id)
			resp, err := call.Do()
			if err != nil {
				fmt.Println("\nERROR FETCHING VID DEETZ: ", err)
			}

			if len(resp.Items) != 0 {
				q.CurrentlyPlayingDetails.Snippet = resp.Items[0].Snippet
				q.CurrentlyPlayingDetails.ContentDetails = resp.Items[0].ContentDetails
				q.CurrentlyPlayingDetails.Statistics = resp.Items[0].Statistics
			} else {
				q.CurrentlyPlayingDetails.Snippet = &youtube.VideoSnippet{
					CategoryId:  "Unknown",
					Title:       "Unknown",
					Description: "Failed to fetch data",
					Thumbnails: &youtube.ThumbnailDetails{
						High: &youtube.Thumbnail{
							Url: "https://i.imgur.com/s36ueeb.jpg",
						},
					},
				}
				q.CurrentlyPlayingDetails.ContentDetails = &youtube.VideoContentDetails{
					Duration: "PT0M0S",
				}
				q.CurrentlyPlayingDetails.Statistics = &youtube.VideoStatistics{
					LikeCount:    uint64(0),
					DislikeCount: uint64(0),
				}
			}

			if len(resp.Items) == 0 {
				fmt.Println("\nNo data retrieved from Youtube|Skipping...")
				q.RemoveQueueEntry()
				wg.Done()
				continue
			}
			es, err := q.GetEncodeSession(requestURL)
			if err != nil {
				fmt.Printf("\nError encoding: %+v\n", err)
				q.RemoveQueueEntry()
				wg.Done()
				continue
			}
			esData, err := es.ReadFrame()
			if err != nil {
				fmt.Println("\nError: ", err)
				q.RemoveQueueEntry()
				wg.Done()
				continue
			}
			if esData == nil {
				fmt.Println("\nNo audio data|Skipping...")
				q.RemoveQueueEntry()
				wg.Done()
				continue
			}
			vc, err = q.establishVoiceConnection(vc, disgordClientAPI, guild, q.VoiceCache[860286976296878080], q.VoiceCache[q.NowPlayingUID])
			if err != nil {
				fmt.Printf("\nError establishing voice connection: %+v\n", err)
				q.RemoveQueueEntry()
				wg.Done()
				continue
			}

			// Ticker needed for smooth opus frame delivery to prevent playback stuttering
			ticker := time.NewTicker(20 * time.Millisecond)
			eofDone := make(chan bool)
			forceDone := make(chan bool)
			eofChannel := make(chan bool)
			stopChannel := make(chan bool)

			// Goroutine just cycles through the opusFrames produced from the encoding process
			// The channels allow for realtime interaction/playback control from events triggered by users
			go func(waitGroup *sync.WaitGroup) {
				fmt.Println("Starting main goRoutine")
				defer es.Cleanup()
				defer ticker.Stop()
				defer waitGroup.Done()
				defer fmt.Println("Leaving main goRoutine")
				for {
					select {
					case <-q.Shuffle:
						ticker.Stop()
						q.stopPlaybackAndTalking(vc, es)
						q.ShuffleNowPlayingQueue()
						return
					case <-q.Stop:
						ticker.Stop()
						q.stopPlaybackAndTalking(vc, es)
						q.EmptyQueue()
						stopChannel <- true
					case <-q.Next:
						ticker.Stop()
						q.stopPlaybackAndTalking(vc, es)
						q.RemoveQueueEntry()
						stopChannel <- true
					case <-eofDone:
						fmt.Print("\neofDone received true, stopping and returning from goRoutine\n")
						q.stopPlaybackAndTalking(vc, es)
						q.RemoveQueueEntry()
						return
					case <-forceDone:
						return
					case channelID := <-q.ChannelHop:
						vc.StopSpeaking()
						vc, err = q.establishVoiceConnection(vc, disgordClientAPI, guild, 0, channelID)
						if err != nil {
							fmt.Printf("\nError establishing voice connection: %+v\n", err)
						}
						err = vc.StartSpeaking()
						if err != nil {
							fmt.Println("Error starting speaking: \n", err)
						}
					case <-q.Pause:
						ticker.Stop()
					case <-q.Play:
						ticker.Reset(20 * time.Millisecond)
						fmt.Println("Resuming...")
					case <-ticker.C:
						nextFrame, err := es.OpusFrame()
						if err != nil && err != io.EOF {
							fmt.Printf("\nError sending next opus frame: %+v\n", err)
						}
						if err == io.EOF {
							fmt.Println("EOF, sending true to eofChannel...")
							eofChannel <- true
							q.stopPlaybackAndTalking(vc, es)
							q.RemoveQueueEntry()
							fmt.Println("Song ended, moving on....")
							return
						}
						vc.SendOpusFrame(nextFrame)
					}
				}
			}(&wg)
			go func(waitGroup *sync.WaitGroup) {
				fmt.Println("Starting secondary goRoutine")
				waitGroup.Add(1)
				defer waitGroup.Done()
				defer fmt.Println("Leaving secondary goRoutine")
				for {
					time.Sleep(1 * time.Second)
					select {
					case <-eofChannel:
						fmt.Print("\neofChannel received true, sending true to eofDone channel\n")
						eofDone <- true
						return
					case <-stopChannel:
						forceDone <- true
						return
					}
				}
			}(&wg)
		}
		wg.Wait()
		if len(q.UserQueue) == 0 {
			q.NowPlayingUID = 0
			q.CurrentlyPlayingDetails = PlayingDetails{}
			q.NowPlayingURL = ""
		}
		// Leave if the bot is the only member in a voice channel
		if len(q.VoiceCache) == 1 && q.VoiceCache[860286976296878080] != 0 && vc != nil {
			vc.Close()
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
		Bitrate:          128,
		Application:      "audio",
		CompressionLevel: 10,
		PacketLoss:       1,
		BufferedFrames:   200,
		VBR:              false,
		StartTime:        0,
		RawOutput:        true,
		Threads:          8,
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
		msgs, err := disgordClient.Channel(1031788884960493618).GetMessages(&disgord.GetMessagesParams{
			Limit: 10,
		})
		if err != nil {
			fmt.Printf("\nCould not get messages from jukebox channel: %+v", err.Error())
		}
		if len(q.UserQueue) > 0 && q.NowPlayingUID != 0 && q.CurrentlyPlayingDetails.Snippet != nil {
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
					1031788884960493618,
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
								{
									Name:  "Requested by",
									Value: requesteeName.Username,
								},
								{
									Name:   "Duration",
									Value:  timeFields[1],
									Inline: true,
								},
								{
									Name:   "Upvotes",
									Value:  likeStr,
									Inline: true,
								},
								{
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
	// This check is in place since the queue may not exist, due to the user leaving voice chat
	// Because when a user leaves voice chat, their queue is deleted
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

// ShuffleNowPlayingQueue reorganizes the order of the queue entries for randomized playback
func (q *Queue) ShuffleNowPlayingQueue() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(q.UserQueue[q.NowPlayingUID]), func(i, j int) {
		q.UserQueue[q.NowPlayingUID][i], q.UserQueue[q.NowPlayingUID][j] = q.UserQueue[q.NowPlayingUID][j], q.UserQueue[q.NowPlayingUID][i]
	})
}

func (q *Queue) ShuffleQueueByID(id disgord.Snowflake) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(q.UserQueue[id]), func(i, j int) {
		q.UserQueue[id][i], q.UserQueue[id][j] = q.UserQueue[id][j], q.UserQueue[id][i]
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
		fmt.Printf("\nError stopping encoding: %+v", err)
	}
	err = vc.StopSpeaking()
	if err != nil && err != io.EOF {
		fmt.Printf("\nError stopping speaking: %+v\n", err)
	}
}

func (q *Queue) establishVoiceConnection(prevVC disgord.VoiceConnection, client disgordiface.DisgordClientAPI, guild disgordiface.GuildQueryBuilderAPI, botChannelID disgord.Snowflake, requesteeChannelID disgord.Snowflake) (disgord.VoiceConnection, error) {
	if botChannelID == 0 {
		vc, err := guild.VoiceChannel(requesteeChannelID).Connect(false, true)
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
		newVC, err := guild.VoiceChannel(requesteeChannelID).Connect(true, false)
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
