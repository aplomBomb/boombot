package discord

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	disgordiface "github.com/aplombomb/boombot/discord/ifaces"
)

// CmdArguments represents the arguments entered by the user after a command
type CmdArguments []string
type msgEvent disgord.Message

// TODO: Get rid of these global variables
var ctx = context.Background()
var disgordGlobalAPI disgordiface.DisgordClientAPI
var disgordGlobalClient *disgord.Client
var globalGuild disgord.GuildQueryBuilder
var ytService *youtube.Service
var session disgord.Session
var globalQueue *Queue

// Version of BoomBot
const Version = "v1.0.0-alpha"

const (
	host   = "db"
	port   = 5432
	dbname = "bomb"
)

func init() {

	fmt.Printf(`
	▄▄▄▄·             • ▌ ▄ ·.  ▄▄▄▄      ▄▄▄▄▄▄▄
	▐█ ▀█▪ ▄█▀▄  ▄█▀▄ ·██ ▐███▪▐█ ▀█▪ ▄█▀▄ •██
	▐█▀▀█▄▐█▌.▐▌▐█▌.▐▌▐█ ▌▐▌▐█·▐█▀▀█▄▐█▌.▐▌ ▐█.▪
	██▄▪▐█▐█▌.▐▌▐█▌.▐▌██ ██▌▐█▌██▄▪▐█▐█▌.▐▌ ▐█▌·
	·▀▀▀▀  ▀█▄▀▪ ▀█▄▀▪▀▀  █▪▀▀▀·▀▀▀▀  ▀█▄▀▪ ▀▀▀ %-16s\/`+"\n\n", Version)
}

// BotRun | Start the bot and react to events
func BotRun(client *disgord.Client, prefix string, gID string, yk string) {
	// dbUser := os.Getenv("POSTGRES_USER")
	// dbPass := os.Getenv("POSTGRES_PASSWORD")
	// pgCreds := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	host, port, dbUser, dbPass, dbname)
	// db, err := sql.Open("postgres", pgCreds)
	// if err != nil {
	// 	log.Fatal("\nError connecting to DB: ", err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	queue := NewQueue(disgord.ParseSnowflakeString(gID))
	globalQueue = queue
	disgordGlobalClient = client
	gg := disgordGlobalClient.Guild(disgord.ParseSnowflakeString(gID))
	globalGuild = gg
	ytService, _ = youtube.NewService(ctx, option.WithAPIKey(yk))
	vlc := ytService.Videos.List([]string{"contentDetails", "snippet", "statistics"})
	filter, _ := std.NewMsgFilter(ctx, client)
	filter.SetPrefix(prefix)
	client.Gateway().WithMiddleware(filter.NotByBot, filter.HasPrefix, std.CopyMsgEvt, filter.StripPrefix).MessageCreate(RespondToCommand)
	client.Gateway().MessageReactionAdd(RespondToReaction)
	client.Gateway().VoiceStateUpdate(RespondToVoiceChannelUpdate)
	client.Gateway().MessageCreate(RespondToMessage)
	// client.Gateway().PresenceUpdate(RespondToPresenceUpdate)
	go globalQueue.ListenAndProcessQueue(ctx, session, client, gg, vlc)
	go globalQueue.ManageJukebox(client)
	defer client.Gateway().StayConnectedUntilInterrupted()
	fmt.Println("BoomBot is running")
}
