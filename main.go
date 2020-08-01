package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aplombomb/boombot/src/pkg/multiplexer"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Version of BoomBot
const Version = "v0.0.0-alpha"

// Session is a global instance of discordgo
// api available for use throughout the app
var Session, _ = discordgo.New()

// Router is global for easy use thoughout the app.
// Passed string will serve as command prefix
var Router = multiplexer.New("**")

func init() {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	//Discord Authentication Token
	Session.Token = os.Getenv("TOKEN")
	if Session.Token == "" {
		flag.StringVar(&Session.Token, "t", "", "Discord Authentication Token")
	}
}

func main() {

	var err error

	// BoomBot cli logo
	fmt.Printf(`
	▄▄▄▄·             • ▌ ▄ ·.  ▄▄▄▄      ▄▄▄▄▄▄▄
	▐█ ▀█▪ ▄█▀▄  ▄█▀▄ ·██ ▐███▪▐█ ▀█▪ ▄█▀▄ •██  
	▐█▀▀█▄▐█▌.▐▌▐█▌.▐▌▐█ ▌▐▌▐█·▐█▀▀█▄▐█▌.▐▌ ▐█.▪
	██▄▪▐█▐█▌.▐▌▐█▌.▐▌██ ██▌▐█▌██▄▪▐█▐█▌.▐▌ ▐█▌·
	·▀▀▀▀  ▀█▄▀▪ ▀█▄▀▪▀▀  █▪▀▀▀·▀▀▀▀  ▀█▄▀▪ ▀▀▀ %-16s\/`+"\n\n", Version)

	// Parse args from command line
	flag.Parse()

	// Check for token
	if Session.Token == "" {
		log.Println("You must provide a Discord authentication token.")
		return
	}

	// Open Discord websocket
	err = Session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanup
	Session.Close()

}
