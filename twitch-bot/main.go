package main

import (
	"encoding/json"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/lpernett/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	cmdFile, fileErr := os.ReadFile("commands.json")
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	var alias = "--"
	var commands map[string]string
	json.Unmarshal(cmdFile, &commands)

	client := twitch.NewClient("tutkuofnight", os.Getenv("PASS"))
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if message.Message[:2] == alias {
			for key, val := range commands {
				if message.Message[2:] == key {
					client.Reply(message.Channel, message.User.ID, val)
				}
			}
		}
	})

	client.Join("tutkuofnight")
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
