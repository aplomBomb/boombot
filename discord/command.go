package discord

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/andersfylling/disgord"

	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
	youtubeiface "github.com/aplombomb/boombot/youtube/ifaces"
)

// CommandEventClient contains the data for all command processing
type CommandEventClient struct {
	data          *disgord.Message
	disgordClient disgordiface.DisgordClientAPI
	youtubeClient youtubeiface.YoutubeSearchServiceAPI
}

// NewCommandEventClient returns a pointer to a new CommandEventClient
func NewCommandEventClient(data *disgord.Message, disgordClient disgordiface.DisgordClientAPI, youtubeClient youtubeiface.YoutubeSearchServiceAPI) *CommandEventClient {
	return &CommandEventClient{
		data:          data,
		disgordClient: disgordClient,
		youtubeClient: youtubeClient,
	}
}

// Delegate evaluates commands and sends them to be processed
func (cec *CommandEventClient) Delegate() {
	cmd, _ := cec.DisectCommand()
	switch cmd {
	case "help", "h", "?", "wtf":
		hcc := NewHelpCommandClient(cec.data, cec.disgordClient)
		hcc.SendHelpMsg()
	case "play":
		// queries := flag.String("query", "deadmau5", "Search term")
		// maxResults := flag.Int64("max-results", 25, "Max YouTube results")
		// call := cec.youtubeClient.List([]string{"id", "snippet"}).Q("deadmau5").MaxResults(*maxResults)
		// response, err := call.Do()
		// if err != nil {
		// 	log.Fatalf("\n\n\nERROR: %+v\n\n\n", err)
		// }
		// fmt.Printf("\n\n\nPAYLOAD: %+v\n\n\n", response)

		// Group video, channel, and playlist results in separate lists.
		// videos := make(map[string]string)
		// channels := make(map[string]string)
		// playlists := make(map[string]string)

		// Iterate through each item and add it to the correct list.
		// for _, item := range response.Items {
		// 	switch item.Id.Kind {
		// 	case "youtube#video":
		// 		videos[item.Id.VideoId] = item.Snippet.Title
		// 	case "youtube#channel":
		// 		channels[item.Id.ChannelId] = item.Snippet.Title
		// 	case "youtube#playlist":
		// 		playlists[item.Id.PlaylistId] = item.Snippet.Title
		// 	}
		// }

		// printIDs("Videos", videos)
		// printIDs("Channels", channels)
		// printIDs("Playlists", playlists)
		resp, err := http.Get("http://localhost:8080/mp3/https://www.youtube.com/watch?v=cF1zJYkBW4A")
		if err != nil {
			log.Fatalf("\n\n\nERROR: %+v\n\n\n", err)
		}
		fmt.Printf("\n\n\nPayload: %+v\n\n\n", resp)
	default:

		uc := NewUnknownCommandClient(cec.data, cec.disgordClient)

		// err := Unknown(data.Message, client)

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
