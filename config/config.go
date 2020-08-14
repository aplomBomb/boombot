package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ConfJSONStruct is used to unmarshal the config.json
type ConfJSONStruct struct {
	Prefix   string `json:"Prefix"`
	BotToken string `json:"Bot_Token"`
	// MongoURL string `json:"Mongo_URL"`
}

// Retrieve reads config from file
func Retrieve(file string) ConfJSONStruct {
	var config ConfJSONStruct
	body, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("error reading config file :", err)
	}

	err = json.Unmarshal(body, &config)
	if err != nil {
		fmt.Println("error unmarshalling config :", err)
	}
	return config
}
