package yt

import (
	youtubeIface "github.com/aplombomb/boombot/youtube/ifaces"
)

// Client represents the collection of data needed to fullfill boombot's youtube functionality
type Client struct {
	YoutubeClient youtubeIface.YoutubeSearchServiceAPI
}

// NewYoutubeClient returns a pointer to a new YtClient
func NewYoutubeClient() *Client {
	return &Client{}
}

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
// resp, err := http.Get("http://localhost:8080/mp3/https://www.youtube.com/watch?v=cF1zJYkBW4A")
// if err != nil {
// 	log.Fatalf("\n\n\nERROR: %+v\n\n\n", err)
// }
// fmt.Printf("\n\n\nPayload: %+v\n\n\n", resp)
