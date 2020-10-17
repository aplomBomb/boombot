package yt

import (
	"context"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// CreateClient returns a Youtube Service client
func CreateClient(token string) (youtube.Service, error) {

	ctx := context.Background()

	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(token))

	if err != nil {
		return youtube.Service{}, err
	}
	return *youtubeService, nil
}
