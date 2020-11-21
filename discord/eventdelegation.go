package discord

import (
	"fmt"
	"strings"
	"time"

	yt "github.com/aplombomb/boombot/youtube"

	"github.com/andersfylling/disgord"
)

// Using this for access to the global clients FOR NOW as passing it through the handlers has proven tricky
// TO-DO find a solution to get rid of the global variables, including the client

// RespondToCommand delegates actions when commands are issued
func RespondToCommand(s disgord.Session, data *disgord.MessageCreate) {
	cec := NewCommandEventClient(data.Message, disgordGlobalClient)
	command, args := cec.DisectCommand()

	fmt.Printf("\nvUserID: %+v\n", data.Message.Author.ID)

	user := data.Message.Author

	fmt.Printf("Command %+v by user %+v | %+v\n", command, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	switch command {
	// TO-DO clean up this god awful repetitive code
	case "play":
		if strings.Contains(args[0], "list=") {
			plis := ytService.PlaylistItems.List([]string{"snippet"})
			ytc := yt.NewYoutubeClient(plis)
			urls, err := ytc.GetPlaylist(args[0])
			if err != nil {
				fmt.Printf("\nERROR GETTING PLAYLIST URLS: %+v\n", err)
			}
			if globalQueue.VoiceCache[data.Message.Author.ID] != 0 {
				resp, err := disgordGlobalClient.SendMsg(
					data.Message.ChannelID,
					&disgord.CreateMessageParams{
						Embed: &disgord.Embed{
							Title:       "**PLAYLIST ACCEPTED**",
							Description: fmt.Sprintf("%+v entries have been added", len(urls)),

							Footer: &disgord.EmbedFooter{
								Text: fmt.Sprintf("*%s's playlist added to queue*", data.Message.Author.Username),
							},
							Timestamp: data.Message.Timestamp,
							Color:     0xeec400,
						},
					},
				)
				if err != nil {
					fmt.Printf("\nERROR ADDING TO QUEUE: %+v\n", err)
				}
				globalQueue.UpdateQueueStateBulk(data.Message.ChannelID, data.Message.Author.ID, urls)
				go deleteMessage(data.Message, 1*time.Second, disgordGlobalClient)
				go deleteMessage(resp, 30*time.Second, disgordGlobalClient)
			} else {
				resp, err := disgordGlobalClient.SendMsg(
					data.Message.ChannelID,
					&disgord.CreateMessageParams{
						Embed: &disgord.Embed{
							Title:       "**PLAYLIST REJECTED**",
							Description: "*You need to be in a voice channel*",

							Footer: &disgord.EmbedFooter{
								Text: fmt.Sprintf("*%s's playlist rejected*", data.Message.Author.Username),
							},
							Timestamp: data.Message.Timestamp,
							Color:     0xeec400,
						},
					},
				)
				if err != nil {
					fmt.Printf("\nERROR ADDING TO QUEUE: %+v\n", err)
				}
				go deleteMessage(data.Message, 1*time.Second, disgordGlobalClient)
				go deleteMessage(resp, 30*time.Second, disgordGlobalClient)
			}

		} else {
			if globalQueue.VoiceCache[data.Message.Author.ID] != 0 {
				resp, err := disgordGlobalClient.SendMsg(
					data.Message.ChannelID,
					&disgord.CreateMessageParams{
						Embed: &disgord.Embed{
							Title:       "**SONG ACCEPTED**",
							Description: "_Song added_",

							Footer: &disgord.EmbedFooter{
								Text: fmt.Sprintf("*%s's song added to queue*", data.Message.Author.Username),
							},
							Timestamp: data.Message.Timestamp,
							Color:     0xeec400,
						},
					},
				)
				if err != nil {
					fmt.Printf("\nERROR ADDING TO QUEUE: %+v\n", err)
				}
				globalQueue.UpdateQueueState(data.Message.ChannelID, data.Message.Author.ID, args[0])
				go deleteMessage(data.Message, 1*time.Second, disgordGlobalClient)
				go deleteMessage(resp, 30*time.Second, disgordGlobalClient)
			} else {
				resp, err := disgordGlobalClient.SendMsg(
					data.Message.ChannelID,
					&disgord.CreateMessageParams{
						Embed: &disgord.Embed{
							Title:       "**SONG REJECTED**",
							Description: "*You need to be in a voice channel*",

							Footer: &disgord.EmbedFooter{
								Text: fmt.Sprintf("*%s's song rejected*", data.Message.Author.Username),
							},
							Timestamp: data.Message.Timestamp,
							Color:     0xeec400,
						},
					},
				)
				if err != nil {
					fmt.Printf("\nERROR ADDING TO QUEUE: %+v\n", err)
				}
				go deleteMessage(data.Message, 1*time.Second, disgordGlobalClient)
				go deleteMessage(resp, 30*time.Second, disgordGlobalClient)
			}

		}
	case "purge":
		resp, err := disgordGlobalClient.SendMsg(
			data.Message.ChannelID,
			&disgord.CreateMessageParams{
				Embed: &disgord.Embed{
					Title:       "**QUEUE PURGED**",
					Description: fmt.Sprintf("%+v entries have been purged", len(globalQueue.UserQueue)),

					Footer: &disgord.EmbedFooter{
						Text: fmt.Sprintf("Purged by %s", data.Message.Author.Username),
					},
					Timestamp: data.Message.Timestamp,
					Color:     0xeec400,
				},
			},
		)
		if err != nil {
			fmt.Printf("\nERROR CREATING PURGE MESSAGE: %+v\n", err)
		}
		go deleteMessage(data.Message, 1*time.Second, disgordGlobalClient)
		go deleteMessage(resp, 30*time.Second, disgordGlobalClient)

		globalQueue.UserQueue = []string{globalQueue.UserQueue[0]}

	default:
		cec.Delegate()
	}

}

// RespondToMessage delegates actions when messages are created
func RespondToMessage(s disgord.Session, data *disgord.MessageCreate) {
	user := data.Message.Author

	fmt.Printf("Message %+v by user %+v | %+v\n", data.Message.Content, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	mec := NewMessageEventClient(data.Message, disgordGlobalClient)
	err := mec.FilterNonModLinks()
	if err != nil {
		fmt.Printf("\nError filtering non-mod link: %+v\n", err)
	}
}

// RespondToReaction delegates actions when reactions are added to messages
func RespondToReaction(s disgord.Session, data *disgord.MessageReactionAdd) {
	userQueryBuilder := disgordGlobalClient.User(data.UserID)

	user, err := userQueryBuilder.Get()

	if err != nil {
		fmt.Printf("\nError getting user: %+v\n", err)
	}
	// fmt.Printf("Message reaction %+v by user %+v | %+v\n", data.PartialEmoji.Name, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	rec := NewReactionEventClient(data.PartialEmoji, data.UserID, data.ChannelID, data.MessageID, disgordGlobalClient)
	msg, err := rec.GenerateModResponse()
	if err != nil {
		fmt.Printf("\nError generating mod reaction response: %+v\n", err)
	}
	//TO-DO Sending the dm here as opposed having it sent via GenerateModResponse for testing purposes
	// Using it here at least allows me to get full coverage of the reactions logic
	// The SendMsg method of disgord.User requires session arg which has proven difficult to mock
	if msg != nil {
		user.SendMsg(ctx, s, msg)
	}
}

// RespondToVoiceChannelUpdate updates the server's voice channel cache every time an update is emitted
func RespondToVoiceChannelUpdate(s disgord.Session, data *disgord.VoiceStateUpdate) {
	globalQueue.UpdateVoiceCache(data.ChannelID, data.UserID)
}
