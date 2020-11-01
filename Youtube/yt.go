package yt

import (
	"context"
	"fmt"
	"time"

	"github.com/andersfylling/disgord"

	"google.golang.org/api/youtube/v3"
)

type YoutubeClient struct {
	service   *youtube.Service
	query     string
	timestamp time.Time
	creator   *disgord.User
}

// CreateClient returns a Youtube Service client
func New(ytS *youtube.Service, query string, creator *disgord.User) (*YoutubeClient, error) {
	return &YoutubeClient{
		service:   ytS,
		query:     query,
		timestamp: time.Now(),
		creator:   creator,
	}, nil
}

func (yt *YoutubeClient) Search(query string) string {
	// result, err := yt.service.Search()
	return "test"
}

func (yt *YoutubeClient) VerifyVoiceChat(sess disgord.Session) bool {

	ctx := context.Background()

	user, err := sess.GetUser(ctx, yt.creator.ID)

	if err != nil {
		fmt.Println("ERROR", err)
	}
	fmt.Printf("\n\nUSER!!: %+v\n\n", user)
	return false
}
