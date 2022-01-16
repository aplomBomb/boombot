package config

import (
	"os"

	"github.com/andersfylling/disgord"
)

// BoombotCreds defines auth structure for boombot
type BoombotCreds struct {
	BotToken     string `json:"BOT_TOKEN"`
	YoutubeToken string `json:"YOUTUBE_TOKEN"`
}

type ServerIDs struct {
	GuildID, JukeboxID, McModChID, TihiID, BattleStationID, VcInitializerID, BoombotID, DiogoID, BugID, RickyID, AverossID, FurryID disgord.Snowflake
}

// GetSecrets retrieves all tokens required by the bot via AWS SecretsManager
func GetSecrets() (*BoombotCreds, error) {
	var sec = BoombotCreds{
		BotToken:     os.Getenv("BOOMBOT_TOKEN"),
		YoutubeToken: os.Getenv("YOUTUBE_TOKEN"),
	}

	return &sec, nil
}

func GetServerIDs() ServerIDs {
	return ServerIDs{
		GuildID:         931430727919231026,
		JukeboxID:       932082717649162291,
		McModChID:       932321289828462652,
		TihiID:          931576984113389638,
		BattleStationID: 932105622483259452,
		VcInitializerID: 931539397600509983,
		BoombotID:       860286976296878080,
		DiogoID:         88482210440544256,
		BugID:           363462448906240001,
		RickyID:         693863839283937373,
		AverossID:       366666843915812874,
		FurryID:         858802176717094984,
	}
}
