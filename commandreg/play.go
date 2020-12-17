package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type CommandInfo struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Options     []CommandOption `json:"options"`
}

type CommandOption struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Type        int             `json:"type"`
	Required    bool            `json:"required"`
	Choices     []CommandChoice `json:"choices"`
}

type CommandChoice struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  int    `json:"type"`
}

func main() {
	botToken := os.Getenv("BOOMBOT_TOKEN")
	guildID := os.Getenv("GUILD_ID")
	// fmt.Println("Token: ", botToken)
	testCommand := CommandInfo{
		Name:        "Test",
		Description: "A test command for BoomBot slash command functionality",
		Options:     []CommandOption{},
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	requestBody, err := json.Marshal(testCommand)
	if err != nil {
		panic(err)
	}

	fmt.Println("Req body: ", string(requestBody))

	reqURL := fmt.Sprintf("https://discord.com/api/v8/applications/739154323015204935/guilds/%+v/commands", guildID)

	request, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(requestBody))
	request.Header.Set("Authorization", fmt.Sprintf("Bot %+v", botToken))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}

	// fmt.Println("Request: ", requestBody)

	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Panicln(string(body))
}
