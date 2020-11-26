package yt

import (
	youtubeiface "github.com/aplombomb/boombot/youtube/ifaces"
)

type VideoClient struct {
	YoutubeVideoClient youtubeiface.YoutubeSearchServiceAPI
}

type VideoClientTest struct {
	YoutubeVideoDetails youtubeiface.YoutubeVideoDetailsAPI
}

func NewVideoClient(vc youtubeiface.YoutubeSearchServiceAPI) *VideoClient {
	return &VideoClient{
		YoutubeVideoClient: vc,
	}
}

func NewVideoClientTest(vcd youtubeiface.YoutubeVideoDetailsAPI) *VideoClientTest {
	return &VideoClientTest{
		YoutubeVideoDetails: vcd,
	}
}

func (c *VideoClient) GetVideoDetails(id string) {
	// resp := c.YoutubeVideoClient.Q(id)
	// slr, err := resp.Do()
	// if err != nil {
	// 	fmt.Println("\nERROR FETCHING SONG INFO: ", err)
	// }
	// fmt.Println("\nSONG NAME: ", slr.Items[0].Snippet.Title)
}
