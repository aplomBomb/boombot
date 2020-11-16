package yt

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	youtubeIface "github.com/aplombomb/boombot/youtube/ifaces"
)

// Client represents the collection of data needed to fullfill boombot's youtube functionality
type Client struct {
	YoutubeClient youtubeIface.YoutubeSearchServiceAPI
}

// NewYoutubeClient returns a pointer to a new YtClient
func NewYoutubeClient(ss youtubeIface.YoutubeSearchServiceAPI) *Client {
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
