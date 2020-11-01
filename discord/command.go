package discord

import (
	"fmt"
	"log"

	yt "github.com/aplombomb/boombot/Youtube"

	"github.com/andersfylling/disgord"
)

//RespondToCommand handles all messages that begin with prefix
func RespondToCommand(s disgord.Session, data *disgord.MessageCreate) {
	cmd, args := ParseMessage(data)

	switch cmd {
	case "help", "h", "?", "wtf":
		fmt.Println(data.Message.Content)
		help(data, args)
	case "play":

		// init the Youtube client
		ytClient, err := yt.New(conf.YoutubeToken, data.Message.Content, data.Message.Author)

		if err != nil {
			log.Fatal("YT API ERROR: ", err)
		}

		fmt.Printf("\nYT Client Created: %+v\n\n\n", ytClient)

		// yt.PrintIt(ytClient)
	default:
		unknown(data)
	}

}
