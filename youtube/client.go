package yt

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	youtubeiface "github.com/aplombomb/boombot/youtube/ifaces"
)

// Client represents the collection of data needed to fullfill boombot's youtube functionality
type Client struct {
	YoutubeClient youtubeiface.YoutubePlaylistItemsServiceAPI
}

// NewYoutubeClient returns a pointer to a new YtClient
func NewYoutubeClient(ss youtubeiface.YoutubePlaylistItemsServiceAPI) *Client {
	return &Client{
		YoutubeClient: ss,
	}
}

// SearchAndDownload returns a string search query based off the provided play command arguments
func (ytc *Client) SearchAndDownload(arg string) (string, error) {
	requestURL := fmt.Sprintf("http://localhost:8080/mp3/%+v", arg)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Fatalf("\n\n\nERROR FETCHING SONG: %+v\n\n\n", err)
	}
	fmt.Printf("\nHEADER: %+v\n", resp.Header)

	filename := fmt.Sprint("song.mp3")

	out, err := os.Create(filename)
	if err != nil {
	}

	io.Copy(out, resp.Body)

	defer out.Close()
	defer resp.Body.Close()

	return filename, nil
}

// GetPlaylist accepts a playlist url and return a slice containing the url's of each video in the playlist
func (ytc *Client) GetPlaylist(arg string) ([]string, error) {
	playlistID := strings.Split(arg, "=")[2]
	fmt.Printf("\nPLAYLISTARG: %+v\n", playlistID)
	songIDs := []string{}
	plic := ytc.YoutubeClient.PlaylistId(playlistID).MaxResults(999)
	resp, err := plic.Do()
	if err != nil {
		return nil, err
	}

	for _, v := range resp.Items {
		url := fmt.Sprintf("https://www.youtube.com/watch?v=%+v", v.Snippet.ResourceId.VideoId)
		songIDs = append(songIDs, url)
	}
	fmt.Printf("\nNUMBEROFSONGS: %+v\n", len(resp.Items))
	fmt.Printf("\nSONGIDS: %+v\n", songIDs)
	// return urls, nil
	return songIDs, nil
}
