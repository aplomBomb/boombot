package discord

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
)

func randomDiogoWarning() string {
	var diogoWarnings = []string{
		"\n**Hey Doogie, have you googled this yet?**",
		"\n**Are you sure this isn't a dumb question, dongo?**",
		"\n**Hey diaygo, have you considered asking someone else?**",
		"\n**Is this a question worth asking?**",
		"\n**Bomb is currently busy, try again never**",
		"\n**Think long and hard about this, deego**",
		"\n**Do you honestly expect a serious answer?**",
		"\n**Maybe.....no?**",
		"\n**Fuck your face dingo**",
		"\n**I'm pretty sure bomb would rather shove shards of glass up his own ass**",
		"\n**If bomb had a dollar for every silly question you've asked, he'd have at least $420**",
		"\n**I believe there's someone a lot dumber that would enjoy answering that instead**",
		"\n**Get lost you filthy degenerate**",
	}

	return diogoWarnings[rand.Intn(len(diogoWarnings))]
}

// Goal is to have this file as small as possible, the purpose of this file isn't about delegation so to say
// It's acting as an intermediary; separating global vars from the business logic
// This means the more code/logic I can get out of this file, the more code/logic that can be tested

// RespondToCommand delegates actions when commands are issued
func respondToCommand(s disgord.Session, data *disgord.MessageCreate) {
	plis := ytService.PlaylistItems.List([]string{"snippet", "status", "contentDetails"})
	ytv := ytService.Search.List([]string{"snippet"}).MaxResults(3).Order("relevance").SafeSearch("none").Type("video")
	cec := NewCommandEventClient(data.Message, disgordGlobalClient, plis, ytv, globalQueue)
	cec.Delegate()
}

// RespondToMessage delegates actions when messages are created
func respondToMessage(s disgord.Session, data *disgord.MessageCreate) {
	// Responses to specific channels
	switch data.Message.ChannelID {
	case ServerIDs.JukeboxID:
		data.Message.React(ctx, s, "\u26D4") // Purge emoji
		time.Sleep(1 * time.Second)
		data.Message.React(ctx, s, "\u267B") // Shuffle Emoji
		time.Sleep(1 * time.Second)
		data.Message.React(ctx, s, "\u23F8") // Pause Emoji
		time.Sleep(1 * time.Second)
		data.Message.React(ctx, s, "\u25B6") // Play Emoji
		time.Sleep(1 * time.Second)
		data.Message.React(ctx, s, "\u23E9") // Next emoji
		time.Sleep(1 * time.Second)
		// data.Message.React(ctx, s, "\u23EC") // Download emoji
		// time.Sleep(1 * time.Second)
	}
	// When diogo has another dumb question for me
	if data.Message.Author.ID == ServerIDs.DiogoID && (strings.Split(data.Message.Content, " ")[0] == "bomb" || strings.Split(data.Message.Content, " ")[1] == "bomb") {
		data.Message.Reply(ctx, s, "<@"+data.Message.Author.ID.String()+">"+randomDiogoWarning())
	}
	if data.Message.Content == "listen here you little shit" {
		data.Message.Reply(ctx, s, "<@"+data.Message.Author.ID.String()+">"+"Here's a big shit https://th.bing.com/th/id/R.a0f1072833b3c8eabee91647b65d227d?rik=k%2boQXoUpVEx0%2fQ&riu=http%3a%2f%2f38.media.tumblr.com%2ftumblr_ll0jboeGDa1qas26so1_500.jpg&ehk=lH80J0kbe7MFLORwwNJ4wquah5gxdkOQfJ%2fGTsjjmOk%3d&risl=&pid=ImgRaw&r=0")
	}

	user := data.Message.Author
	fmt.Printf("Message %+v by %+v | %+v\n", data.Message.Content, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	mec := NewMessageEventClient(data.Message, disgordGlobalClient)
	err := mec.FilterMessages()
	if err != nil {
		fmt.Printf("\nError filtering non-mod link: %+v\n", err)
	}
}

// RespondToReaction delegates actions when reactions are added to messages
func respondToReaction(s disgord.Session, data *disgord.MessageReactionAdd) {
	userQueryBuilder := disgordGlobalClient.User(data.UserID)
	user, err := userQueryBuilder.Get()
	if err != nil {
		fmt.Printf("\nError getting user: %+v\n", err)
	}
	// fmt.Printf("Message reaction %+v by user %+v | %+v\n", data.PartialEmoji.Name, user.Username, time.Now().Format("Mon Jan _2 15:04:05 2006"))
	rec := NewReactionEventClient(data.PartialEmoji, data.UserID, data.ChannelID, data.MessageID, disgordGlobalClient)
	rec.HandleJukeboxReact(globalQueue)
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
func respondToVoiceChannelUpdate(s disgord.Session, data *disgord.VoiceStateUpdate) {
	globalQueue.UpdateVoiceCache(data.ChannelID, data.UserID)
	if data.ChannelID != 0 && data.ChannelID != globalQueue.VoiceCache[ServerIDs.BoombotID] && data.UserID == globalQueue.NowPlayingUID && globalQueue.VoiceCache[data.UserID] != 0 {
		go func() {
			globalQueue.ChannelHop <- data.ChannelID
		}()
	}
}

// RespondToPresenceUpdate fires when a server member's presence state changes
func respondToPresenceUpdate(s disgord.Session, data *disgord.PresenceUpdate) {

	activity, err := data.Game()

	if err != nil {
		fmt.Printf("\nError when feting presence update for %+v | \nError: %+v\n", data.User.Username, err.Error())
	}

	fmt.Printf("\n%+v\n", data.User.Username)
	fmt.Printf("\nActivity Name: %+v\n", activity.Name)
	fmt.Printf("\nActivity AppID: %+v\n", activity.ApplicationID)
	fmt.Printf("\nActivity Details: %+v\n", activity.Details)
	fmt.Printf("\nActivity Party State: %+v\n", activity.State)
	fmt.Printf("\nActivity Emoji: %+v\n", activity.Emoji)
	fmt.Printf("\nActivity Type: %+v\n", activity.Type)
}
