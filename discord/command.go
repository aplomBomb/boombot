package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
	yt "github.com/aplombomb/boombot/youtube"
	youtubeiface "github.com/aplombomb/boombot/youtube/ifaces"
)

// CommandEventClient contains the data for all command processing
type CommandEventClient struct {
	data          *disgord.Message
	disgordClient disgordiface.DisgordClientAPI
	ytps          youtubeiface.YoutubePlaylistServiceAPI
	ytss          youtubeiface.YoutubeSearchAPI
	queue         disgordiface.QueueClientAPI
}

// NewCommandEventClient returns a pointer to a new CommandEventClient
func NewCommandEventClient(data *disgord.Message, disgordClient disgordiface.DisgordClientAPI, ytps youtubeiface.YoutubePlaylistServiceAPI, ytss youtubeiface.YoutubeSearchAPI, qc disgordiface.QueueClientAPI) *CommandEventClient {
	return &CommandEventClient{
		data:          data,
		disgordClient: disgordClient,
		ytps:          ytps,
		ytss:          ytss,
		queue:         qc,
	}
}

// Delegate evaluates commands and sends them to be processed
func (cec *CommandEventClient) Delegate() {
	mec := NewMessageEventClient(cec.data, cec.disgordClient)
	cmd, args := cec.DisectCommand()
	switch cmd {
	case "help", "h", "?", "wtf":
		hcc := NewHelpCommandClient(cec.data, cec.disgordClient)
		hcc.SendHelpMsg()
	case "next":
		cec.queue.TriggerNext()
		go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
	case "shuffle":
		cec.queue.TriggerShuffle()
		go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
	case "stop":
		cec.queue.TriggerStop()
		go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
	case "play":
		// CHECK IF USER IS IN VOICE CHAT
		if cec.queue.ReturnVoiceCacheEntry(cec.data.Author.ID) == 0 {
			_, err := mec.SendEmbedMsgReply(disgord.Embed{
				Title:       "**Request Rejected**",
				Description: "You need to be in a voice channel to make a request",
				Timestamp:   cec.data.Timestamp,
				Footer: &disgord.EmbedFooter{
					Text: fmt.Sprintf("%+v is not in a voice channel", cec.data.Author.Username),
				},
				Color: 0xeec400,
			},
			)
			if err != nil {
				fmt.Printf("\nError sending request rejected message: %+v\n", err)
			}
			go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
			return
		}
		// CHECK IF USER SUPPLIED ARGUMENT
		if len(args) == 0 {
			_, err := mec.SendEmbedMsgReply(disgord.Embed{
				Title:       "**Empty Request**",
				Description: "You used the play command, but didn't provide an argument",
				Timestamp:   disgord.Time{Time: time.Now()},
				Footer: &disgord.EmbedFooter{
					Text: fmt.Sprintf("%+v only provided the play command", cec.data.Author.Username),
				},
				Color: 0xeec400,
			})
			if err != nil {
				fmt.Printf("\nError sending embed message: %+v\n", err)
			}
			go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
			return
		}
		// PARSE ARGUMENT
		parsedArgs, isURL, err := cec.ParseYoutubeArg(args)
		if err != nil {
			fmt.Println("\nError parsing play argument: ", err)
		}
		switch isURL {
		case true:
			fmt.Println("\nArgs: ", parsedArgs)
			// Check for youtu.be links, extract the ID and append to standard URL for queue processing
			if strings.Contains(parsedArgs[0], "youtu.be") {
				if cec.queue.ReturnVoiceCacheEntry(cec.data.Author.ID) != 0 {
					fields := strings.Split(parsedArgs[0], "youtu.be/")
					// Format the argument to work with the queue's url processor
					requestURL := fmt.Sprintf("https://www.youtube.com/watch?v=%+v", fields[1])
					_, err := mec.SendEmbedMsgReply(disgord.Embed{
						Title:       "**Request Accepted**",
						Description: "Your request has been submitted and will be played soon",
						Timestamp:   cec.data.Timestamp,
						Footer: &disgord.EmbedFooter{
							Text: fmt.Sprintf("%+v added a song to their queue", cec.data.Author.Username),
						},
						Color: 0xeec400,
					},
					)
					if err != nil {
						fmt.Printf("\nError sending request accepted message: %+v\n", err)
					}
					cec.queue.UpdateUserQueueState(cec.data.ChannelID, cec.data.Author.ID, requestURL)
					go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
					return
				}
			}
			// Check if URL is a playlist
			if strings.Contains(parsedArgs[0], "list=") {
				ytplc := yt.NewYoutubePlaylistClient(cec.ytps)
				urls, err := ytplc.GetPlaylist(parsedArgs[0])
				if err != nil {
					fmt.Printf("\nError getting playlist URLs: %+v\n", err)
				}
				if len(urls) != 0 {
					_, err := mec.SendEmbedMsgReply(disgord.Embed{
						Title:       "**Playlist Accepted**",
						Description: fmt.Sprintf("%+v entries have been added", len(urls)),
						Footer: &disgord.EmbedFooter{
							Text: fmt.Sprintf("*%s's playlist added to queue*", cec.data.Author.Username),
						},
						Timestamp: cec.data.Timestamp,
						Color:     0xeec400,
					},
					)
					if err != nil {
						fmt.Printf("\nError: %+v\n", err)
					}
					cec.queue.UpdateUserQueueStateBulk(cec.data.ChannelID, cec.data.Author.ID, urls)
					go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
					return
				}

			} else {
				if cec.queue.ReturnVoiceCacheEntry(cec.data.Author.ID) != 0 {
					_, err := mec.SendEmbedMsgReply(disgord.Embed{
						Title:       "**Request Accepted**",
						Description: "_Song added_",
						Footer: &disgord.EmbedFooter{
							Text: fmt.Sprintf("*%s's request added to queue*", cec.data.Author.Username),
						},
						Timestamp: cec.data.Timestamp,
						Color:     0xeec400,
					},
					)
					if err != nil {
						fmt.Printf("\nError sending song accepted message: %+v\n", err)
					}
					cec.queue.UpdateUserQueueState(cec.data.ChannelID, cec.data.Author.ID, parsedArgs[0])
					go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
					return
				}
			}
		case false:
			fmt.Sprintf("\n==============================\n%+v searched for '%+v'\n==============================", cec.data.Author.Username, cec.data.Content)
			argString := strings.Join(parsedArgs, " ")
			slc := cec.ytss.Q(argString)
			resp, err := slc.Do()
			if err != nil {
				fmt.Println("\nError searching, ", err)
			}
			if len(resp.Items) != 0 {
				fmt.Println("\nITEMS: ", resp.Items)
				vidID := resp.Items[0].Id.VideoId
				url := fmt.Sprintf("\nhttps://www.youtube.com/watch?v=%+v", vidID)
				cec.queue.UpdateUserQueueState(cec.data.ChannelID, cec.data.Author.ID, url)
				fmt.Println("\nURL from search: ", resp.Items[0].Snippet.Title)
				avatarURL, err := cec.data.Author.AvatarURL(64, true)
				if err != nil {
					fmt.Println("\n", err)
				}
				_, err = mec.SendEmbedMsgReply(disgord.Embed{
					Title: resp.Items[0].Snippet.Title,
					Thumbnail: &disgord.EmbedThumbnail{
						URL:    avatarURL,
						Height: 64,
						Width:  64,
					},
					Image: &disgord.EmbedImage{
						URL:    resp.Items[0].Snippet.Thumbnails.Default.Url,
						Height: 128,
						Width:  128,
					},
					Footer: &disgord.EmbedFooter{
						Text: fmt.Sprintf("Added by %s", cec.data.Author.Username),
					},
					Timestamp: cec.data.Timestamp,
					Color:     0xeec400,
				},
				)
				if err != nil {
					fmt.Printf("\nError sending searched song message: %+v\n", err)
				}
				go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
				return
			}
			_, err = mec.SendEmbedMsgReply(disgord.Embed{
				Title:       "**No Results**",
				Description: "No results found",
				Footer: &disgord.EmbedFooter{
					Text: fmt.Sprintf("%s's taste in music is too exotic", cec.data.Author.Username),
				},
				Timestamp: cec.data.Timestamp,
				Color:     0xeec400,
			},
			)
			if err != nil {
				fmt.Printf("\nError sending no search results found messageg: %+v\n", err)
			}
			go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
			return
		}

	case "purge":
		uq := cec.queue.ReturnUserQueue()
		_, err := mec.SendEmbedMsgReply(disgord.Embed{
			Title:       "**Queue Purged**",
			Description: fmt.Sprintf("%+v entries have been purged", len(uq)),

			Footer: &disgord.EmbedFooter{
				Text: fmt.Sprintf("Purged by %s", cec.data.Author.Username),
			},
			Timestamp: cec.data.Timestamp,
			Color:     0xeec400,
		},
		)
		if err != nil {
			fmt.Printf("\nError sending purge message: %+v\n", err)
		}
		go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
		delete(cec.queue.ReturnUserQueue(), cec.queue.ReturnNowPlayingID())
	default:
		uc := NewUnknownCommandClient(cec.data, cec.disgordClient)
		uc.RespondToChannel()
	}
}

// DisectCommand returns the used command and all extraneous arguments
func (cec *CommandEventClient) DisectCommand() (string, []string) {
	var command string
	var args []string
	if len(cec.data.Content) > 0 {
		command = strings.ToLower(strings.Fields(cec.data.Content)[0])
		if len(cec.data.Content) > 1 {
			args = strings.Fields(cec.data.Content)[1:]
		}
	}
	return command, args
}

// ParseYoutubeArg handles argument parsing for the play command
func (cec *CommandEventClient) ParseYoutubeArg(args []string) ([]string, bool, error) {
	parsedArg := []string{}
	isURL := false
	if strings.Contains(args[0], "https://www.youtu") != false {
		isURL = true
		parsedArg = args
		return parsedArg, isURL, nil
	}
	parsedArg = args
	return parsedArg, isURL, nil
}

// Print the ID and title of each result in a list as well as a name that
// identifies the list. For example, print the word section name "Videos"
// above a list of video search results, followed by the video ID and title
// of each matching video.
func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}
