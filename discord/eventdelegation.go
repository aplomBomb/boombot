package discord

import (
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
)

// Goal is to have this file as small as possible, the purpose of this file isn't about delegation so to say
// It's acting as an intermediary; separating global vars from the business logic
// This means the more code/logic I can get out of this file, the more code/logic that can be unit tested via dependency injection

// RespondToCommand delegates actions when commands are issued
func RespondToCommand(s disgord.Session, data *disgord.MessageCreate) {
	plis := ytService.PlaylistItems.List([]string{"snippet", "status", "contentDetails"})
	ytv := ytService.Search.List([]string{"snippet"}).MaxResults(3).Order("viewCount").SafeSearch("none").Type("video")
	cec := NewCommandEventClient(data.Message, disgordGlobalClient, plis, ytv, globalQueue)
	cec.Delegate()
}

// RespondToMessage delegates actions when messages are created
func RespondToMessage(s disgord.Session, data *disgord.MessageCreate) {
	switch data.Message.ChannelID {
	case 779836590503624734:
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
	}
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
func RespondToVoiceChannelUpdate(s disgord.Session, data *disgord.VoiceStateUpdate) {
	globalQueue.UpdateVoiceCache(data.ChannelID, data.UserID)
	if data.ChannelID != 0 && data.ChannelID != globalQueue.VoiceCache[739154323015204935] && data.UserID == globalQueue.NowPlayingUID && globalQueue.VoiceCache[data.UserID] != 0 {
		go func() {
			globalQueue.ChannelHop <- data.ChannelID
			return
		}()
	}
}

// RespondToPresenceUpdate fires when a server member's presence state changes
func RespondToPresenceUpdate(s disgord.Session, data *disgord.PresenceUpdate) {
	// This is temporary for testing\will be moved to a "roles" client
	// after i decide how i want to handle these events

	// Roles drg, twerkov, crafter respectively
	// Will map game name as string key/snowflake of role as value
	// Will make it really easy to remove/add role ids on events
	// roleCache := map[string]disgord.Snowflake

	managedRoles := []disgord.Snowflake{787758251574820864, 737467990647373827, 735890320348282880}
	gameEvent, _ := data.Game()
	fmt.Println("GameEvent: ", gameEvent)
	if gameEvent == nil {
		fmt.Println("This must have been an online/offline event")
		return
	}
	fmt.Println("Game Name: ", gameEvent.Name)
	if len(data.Activities) == 0 {
		memberQueryBuilder := globalGuild.Member(data.User.ID)
		member, err := globalGuild.Member(data.User.ID).Get()
		if err != nil {
			fmt.Println("Error fetching member for role adjustment: ", err)
		}
		roles := member.Roles
		fmt.Printf("\n%+v's roles before: %+v", data.User.Username, roles)
		for _, rv := range roles {
			for _, mrv := range managedRoles {
				if mrv == rv {
					memberQueryBuilder.RemoveRole(mrv)
				}
			}
		}
		fmt.Printf("\n%+v's roles after: %+v", data.User.Username, roles)
		return
	}

	// userID := data.User.ID
	// activityName := data.Activities[0].Name
	// drgRoleID := 787758251574820864

	// for k := range data.Activities {
	// 	fmt.Println("activity: ", data.Activities[k])
	// }
}
