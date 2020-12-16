package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type CommandInfo struct {
	Name        string          `json:"name"`
	Description string          `json:"Description"`
	Options     []CommandOption `json:"Options"`
}

type CommandOption struct {
	Name        string          `json:"Name"`
	Description string          `json:"Description"`
	Type        int             `json:"Type"`
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
	// fmt.Println("Token: ", botToken)
	testCommand := CommandInfo{
		Name:        "Test",
		Description: "A test command for BoomBot's slash command functionality",
		Options: []CommandOption{
			CommandOption{
				Name:        "Option1",
				Description: "First option",
				Type:        2,
				Required:    true,
				Choices: []CommandChoice{
					CommandChoice{
						Name:  "First Choice",
						Value: "first_choice",
						Type:  3,
					},
					CommandChoice{
						Name:  "Second Choice",
						Value: "second_choice",
						Type:  1,
					},
				},
			},
		},
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

	request, err := http.NewRequest("POST", "https://discord.com/api/v8/applications/739154323015204935/commands", bytes.NewBuffer(requestBody))
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Panicln(string(body))
}
