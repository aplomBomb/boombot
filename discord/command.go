package discord

import (
	"context"
	"fmt"
	"log"

	yt "github.com/aplombomb/boombot/Youtube"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/andersfylling/disgord"
)

//RespondToCommand handles all messages that begin with prefix
func RespondToCommand(s disgord.Session, data *disgord.MessageCreate) {
	cmd, args := ParseMessage(data)

	switch cmd {
	case "help", "h", "?", "wtf":
		fmt.Println(data.Message.Content)
		hcc := NewHelpCommandClient(data, client)

		hcc.SendHelpMsg(args)
	case "play":

		// init the Youtube client here for test coverage's sake | will find another home for this later
		ctx := context.Background()

		youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(conf.YoutubeToken))

		if err != nil {
			fmt.Println(err)
		}
		ytClient, err := yt.New(youtubeService, data.Message.Content, data.Message.Author)

		if err != nil {
			log.Fatal("YT API ERROR: ", err)
		}

		fmt.Printf("\nYT Client Created: %+v\n\n\n", ytClient)

		// if inVoice := ytClient.VerifyVoiceChat(s); inVoice == false {
		// 	fmt.Println("User is not in voice channel")
		// } else {
		// 	fmt.Println("User is in voice channel")
		// }

		// yt.PrintIt(ytClient)
	default:

		uc := NewUnknownCommandClient(data.Message, client)

		// err := Unknown(data.Message, client)

		uc.RespondToChannel()
	}

}
