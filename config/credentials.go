package config

import "os"

// BoombotCreds defines auth structure for boombot
type BoombotCreds struct {
	BotToken     string `json:"BOT_TOKEN"`
	YoutubeToken string `json:"YOUTUBE_TOKEN"`
}

// GetSecrets retrieves all tokens required by the bot via AWS SecretsManager
func GetSecrets() (*BoombotCreds, error) {
	var sec = BoombotCreds{
		BotToken:     os.Getenv("BOOMBOT_TOKEN"),
		YoutubeToken: os.Getenv("YOUTUBE_TOKEN"),
	}

	return &sec, nil
}
