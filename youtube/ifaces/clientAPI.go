package yt

import (
	"google.golang.org/api/youtube/v3"
)

// YoutubeServiceAPI provides an interface for mocking youtube service behavior
// type YoutubeServiceAPI interface {
// 	YoutubeVideoDetailsAPI
// 	YoutubeVideoItemServiceAPI
// 	YoutubePlaylistItemsServiceAPI
// }

// YoutubeSearchAPI provides an interface for youtube search behavior
type YoutubeSearchAPI interface {
	Q(q string) *youtube.SearchListCall
}

// YoutubeVideoDetailsAPI provides an interface for mocking
type YoutubeVideoDetailsAPI interface {
	Q(q string) *youtube.VideosListCall
}

// YoutubePlaylistServiceAPI provides an interface for mocking youtube playlistitemscall behavior
type YoutubePlaylistServiceAPI interface {
	// List(part []string) *youtube.PlaylistItemsListCall
	PlaylistId(playlistID string) *youtube.PlaylistItemsListCall
}

// YoutubeVideoServiceAPI provides an interface for mocking youtube videoslistcall behavior
type YoutubeVideoServiceAPI interface {
	List(part []string) *youtube.VideosListCall
	Q(q string) *youtube.VideosListCall
}
