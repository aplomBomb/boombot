package yt

import (
	"context"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// YoutubeClientAPI provides an interface for mocking youtube service behavior
type YoutubeClientAPI interface {
	Search(ctx context.Context, opts ...option.ClientOption) (*youtube.Service, error)
}
