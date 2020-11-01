package yt

import (
	"context"
	"time"

	"github.com/andersfylling/disgord"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeClient struct {
	service   *youtube.Service
	query     string
	timestamp time.Time
	creator   *disgord.User
}

// CreateClient returns a Youtube Service client
func New(token string, query string, creator *disgord.User) (*YoutubeClient, error) {

	ctx := context.Background()

	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(token))

	if err != nil {
		return nil, err
	}

	return &YoutubeClient{
		service:   youtubeService,
		query:     query,
		timestamp: time.Now(),
		creator:   creator,
	}, nil
}

func (yt *YoutubeClient) Search(query string) string {
	// result, err := yt.service.Search()
	return "test"
}
