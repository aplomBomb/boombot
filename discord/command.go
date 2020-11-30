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
	data                         *disgord.Message
	disgordClient                disgordiface.DisgordClientAPI
	youtubePlaylistServiceClient youtubeiface.YoutubePlaylistServiceAPI
	// youtubeVideoServiceClient    youtubeiface.YoutubeVideoServiceAPI
	queue disgordiface.QueueClientAPI
}

// NewCommandEventClient returns a pointer to a new CommandEventClient
func NewCommandEventClient(data *disgord.Message, disgordClient disgordiface.DisgordClientAPI, ytps youtubeiface.YoutubePlaylistServiceAPI, qc disgordiface.QueueClientAPI) *CommandEventClient {
	return &CommandEventClient{
		data:                         data,
		disgordClient:                disgordClient,
		youtubePlaylistServiceClient: ytps,
		// youtubeVideoServiceClient:    ytvs,
		queue: qc,
	}
}

// Delegate evaluates commands and sends them to be processed
func (cec *CommandEventClient) Delegate() {
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
		if len(args) == 0 {
			resp, err := cec.disgordClient.SendMsg(
				cec.data.ChannelID,
				&disgord.CreateMessageParams{
					Embed: &disgord.Embed{
						Title:       "**Empty Request**",
						Description: "*You didn't request anything*",

						Footer: &disgord.EmbedFooter{
							Text: fmt.Sprintf("*%+v only provided the play command*", cec.data.Author.Username),
						},
						Timestamp: cec.data.Timestamp,
						Color:     0xeec400,
					},
				},
			)
			if err != nil {
				fmt.Printf("\nERROR ADDING TO QUEUE: %+v\n", err)
			}
			go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
			go deleteMessage(resp, 7*time.Second, cec.disgordClient)
			return
		}
		if strings.Contains(args[0], "youtu.be") {

			if cec.queue.ReturnVoiceCacheEntry(cec.data.Author.ID) != 0 {
				fields := strings.Split(args[0], "youtu.be/")
				requestURL := fmt.Sprintf("https://www.youtube.com/watch?v=%+v", fields[1])
				fmt.Println("\nRequest: ", requestURL)
				resp, err := cec.disgordClient.SendMsg(
					cec.data.ChannelID,
					&disgord.CreateMessageParams{
						Embed: &disgord.Embed{
							Title:       "**Song Accepted**",
							Description: "*Song added to your queue*",

							Footer: &disgord.EmbedFooter{
								Text: fmt.Sprintf("%+v added a song to their queue", cec.data.Author.Username),
							},
							Timestamp: cec.data.Timestamp,
							Color:     0xeec400,
						},
					},
				)
				if err != nil {
					fmt.Printf("\nERROR ADDING TO QUEUE: %+v\n", err)
				}
				cec.queue.UpdateUserQueueState(cec.data.ChannelID, cec.data.Author.ID, requestURL)
				go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
				go deleteMessage(resp, 7*time.Second, cec.disgordClient)
				return
			} else {
				resp, err := cec.disgordClient.SendMsg(
					cec.data.ChannelID,
					&disgord.CreateMessageParams{
						Embed: &disgord.Embed{
							Title:       "**Song Rejected**",
							Description: "*You need to be in a voice channel*",

							Footer: &disgord.EmbedFooter{
								Text: fmt.Sprintf("*%s's song rejected*", cec.data.Author.Username),
							},
							Timestamp: cec.data.Timestamp,
							Color:     0xeec400,
						},
					},
				)
				if err != nil {
					fmt.Printf("\nERROR ADDING TO QUEUE: %+v\n", err)
				}
				go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
				go deleteMessage(resp, 7*time.Second, cec.disgordClient)
			}
		}
		if strings.Contains(args[0], "list=") {

			ytplc := yt.NewYoutubePlaylistClient(cec.youtubePlaylistServiceClient)

			urls, err := ytplc.GetPlaylist(args[0])
			if err != nil {
				fmt.Printf("\nERROR GETTING PLAYLIST URLS: %+v\n", err)
			}
			if cec.queue.ReturnVoiceCacheEntry(cec.data.Author.ID) != 0 {
				resp, err := cec.disgordClient.SendMsg(
					cec.data.ChannelID,
					&disgord.CreateMessageParams{
						Embed: &disgord.Embed{
							Title:       "**Playlist Accepted**",
							Description: fmt.Sprintf("%+v entries have been added", len(urls)),

							Footer: &disgord.EmbedFooter{
								Text: fmt.Sprintf("*%s's playlist added to queue*", cec.data.Author.Username),
							},
							Timestamp: cec.data.Timestamp,
							Color:     0xeec400,
						},
					},
				)
				if err != nil {
					fmt.Printf("\nERROR ADDING TO QUEUE: %+v\n", err)
				}
				cec.queue.UpdateUserQueueStateBulk(cec.data.ChannelID, cec.data.Author.ID, urls)
				go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
				go deleteMessage(resp, 7*time.Second, cec.disgordClient)
			} else {
				resp, err := cec.disgordClient.SendMsg(
					cec.data.ChannelID,
					&disgord.CreateMessageParams{
						Embed: &disgord.Embed{
							Title:       "**Playlist Rejected**",
							Description: "*You need to be in a voice channel*",

							Footer: &disgord.EmbedFooter{
								Text: fmt.Sprintf("*%s's playlist rejected*", cec.data.Author.Username),
							},
							Timestamp: cec.data.Timestamp,
							Color:     0xeec400,
						},
					},
				)
				if err != nil {
					fmt.Printf("\nERROR ADDING TO QUEUE: %+v\n", err)
				}
				go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
				go deleteMessage(resp, 7*time.Second, cec.disgordClient)
			}
		} else {
			if cec.queue.ReturnVoiceCacheEntry(cec.data.Author.ID) != 0 {
				resp, err := cec.disgordClient.SendMsg(
					cec.data.ChannelID,
					&disgord.CreateMessageParams{
						Embed: &disgord.Embed{
							Title:       "**Song Accepted**",
							Description: "_Song added_",

							Footer: &disgord.EmbedFooter{
								Text: fmt.Sprintf("*%s's song added to queue*", cec.data.Author.Username),
							},
							Timestamp: cec.data.Timestamp,
							Color:     0xeec400,
						},
					},
				)
				if err != nil {
					fmt.Printf("\nERROR ADDING TO QUEUE: %+v\n", err)
				}
				cec.queue.UpdateUserQueueState(cec.data.ChannelID, cec.data.Author.ID, args[0])
				go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
				go deleteMessage(resp, 7*time.Second, cec.disgordClient)
			} else {
				resp, err := cec.disgordClient.SendMsg(
					cec.data.ChannelID,
					&disgord.CreateMessageParams{
						Embed: &disgord.Embed{
							Title:       "**Song Rejected**",
							Description: "*You need to be in a voice channel*",

							Footer: &disgord.EmbedFooter{
								Text: fmt.Sprintf("*%s's song rejected*", cec.data.Author.Username),
							},
							Timestamp: cec.data.Timestamp,
							Color:     0xeec400,
						},
					},
				)
				if err != nil {
					fmt.Printf("\nERROR ADDING TO QUEUE: %+v\n", err)
				}
				go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
				go deleteMessage(resp, 30*time.Second, cec.disgordClient)
			}
		}
	case "purge":
		uq := cec.queue.ReturnUserQueue()
		resp, err := cec.disgordClient.SendMsg(
			cec.data.ChannelID,
			&disgord.CreateMessageParams{
				Embed: &disgord.Embed{
					Title:       "**Queue Purged**",
					Description: fmt.Sprintf("%+v entries have been purged", len(uq)-1),

					Footer: &disgord.EmbedFooter{
						Text: fmt.Sprintf("Purged by %s", cec.data.Author.Username),
					},
					Timestamp: cec.data.Timestamp,
					Color:     0xeec400,
				},
			},
		)
		if err != nil {
			fmt.Printf("\nERROR CREATING PURGE MESSAGE: %+v\n", err)
		}
		go deleteMessage(cec.data, 1*time.Second, cec.disgordClient)
		go deleteMessage(resp, 7*time.Second, cec.disgordClient)
		// cec.queue.UserQueue[cec.queue.NowPlayingUID] = []string{cec.queue.UserQueue[cec.queue.NowPlayingUID][0]}
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
