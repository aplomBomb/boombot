package yt

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	youtubeiface "github.com/aplombomb/boombot/youtube/ifaces"
)

// PlaylistClient represents the collection of data needed to fullfill boombot's youtube functionality
type PlaylistClient struct {
	YoutubePlaylistClient youtubeiface.YoutubePlaylistServiceAPI
}

// NewYoutubePlaylistClient returns a pointer to a new YtClient
func NewYoutubePlaylistClient(pls youtubeiface.YoutubePlaylistServiceAPI) *PlaylistClient {
	return &PlaylistClient{
		YoutubePlaylistClient: pls,
	}
}

// SearchAndDownload returns a string search query based off the provided play command arguments
func (c *PlaylistClient) SearchAndDownload(arg string) (string, error) {
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

// GetPlaylist accepts a playlist url and returns a
// Slice containing the url's of each video in the playlist
func (c *PlaylistClient) GetPlaylist(arg string) ([]string, error) {
	songIndex := 0
	playlistID := ""
	urlFields := strings.Split(arg, "=")
	fmt.Printf("\nPLAYLISTFIELDS: %+v\n", urlFields)
	fmt.Printf("\nPLAYLISTARGLENGTH: %+v\n", len(urlFields))
	// This is kind of disgusting/clean up at some point
	switch len(urlFields) {
	case 2:
		fmt.Printf("\nPLAYLISTID: %+v\n", urlFields[1])
		playlistID = urlFields[1]
	case 3:
		fmt.Printf("\nPLAYLISTID: %+v\n", urlFields[2])
		playlistID = urlFields[2]
	case 4:
		if strings.Contains(urlFields[2], "&index") {
			id := strings.Split(urlFields[2], "&index")
			playlistID = id[0]
			songIndex, _ = strconv.Atoi(urlFields[3])
		}
		if strings.Contains(urlFields[2], "&start_radio") {
			id := strings.Split(urlFields[2], "&start_radio")
			playlistID = id[0]
			songIndex, _ = strconv.Atoi(urlFields[3])
		}
	case 5:
		if strings.Contains(urlFields[3], "ab_channel") {
			urlFields = urlFields[:len(urlFields)-2]
		}
	}
	URLS, err := c.aggregateIDS(playlistID)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\nNUMBEROFSONGSFETCHED: %+v\n", len(URLS))

	// Remove the number of track URLS from the slice to start playlist at the requested index
	if songIndex > 0 {
		for i := 0; i < songIndex; i++ {
			copy(URLS[0:], URLS[0+1:])
			URLS[len(URLS)-1] = ""
			URLS = URLS[:len(URLS)-1]
		}
	}
	for i, v := range URLS {
		fmt.Printf("\n%+v: %+v", i, v)
	}
	return URLS, nil
}

// TO-DO refactor this to fetch track concurrently for fast grabbing of large playlists
func (c *PlaylistClient) aggregateIDS(plID string) ([]string, error) {
	IDS := []string{}
	nextPageToken := ""
	for {
		plic := c.YoutubePlaylistClient.PlaylistId(plID).MaxResults(50).PageToken(nextPageToken)
		resp, err := plic.Do()
		if err != nil {
			return nil, err
		}
		for _, v := range resp.Items {
			// Gonna handle this in the 'GetVideoDetails' method
			if v.Status.PrivacyStatus == "private" || v.Status.PrivacyStatus == "unlisted" {
				continue
			}
			url := fmt.Sprintf("https://www.youtube.com/watch?v=%+v", v.Snippet.ResourceId.VideoId)
			if url != fmt.Sprintf("https://www.youtube.com/watch?v=%+v", "") {
				IDS = append(IDS, url)
			}
		}
		nextPageToken = resp.NextPageToken
		if nextPageToken == "" {
			break
		}
		fmt.Println(".................")
	}
	return IDS, nil
}
