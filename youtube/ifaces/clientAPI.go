package yt

import (
	"google.golang.org/api/youtube/v3"
)

// YoutubeSearchServiceAPI provides an interface for mocking youtube search service service behavior
type YoutubeSearchServiceAPI interface {
	List(part []string) *youtube.SearchListCall
}

// YoutubePlaylistItemsServiceAPI provides and interface for mocking youtube playlistitems service behavior
type YoutubePlaylistItemsServiceAPI interface {
	PlaylistId(playlistID string) *youtube.PlaylistItemsListCall
	// Do(opts ...googleapi.CallOption) (*youtube.PlaylistItemListResponse, error)
}
