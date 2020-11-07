package yt

import (
	"google.golang.org/api/youtube/v3"
)

// YoutubeSearchServiceAPI provides an interface for mocking youtube search service service behavior
type YoutubeSearchServiceAPI interface {
	List(part []string) *youtube.SearchListCall
}
