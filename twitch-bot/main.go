package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/lpernett/godotenv"
	"log"
	"os"
)

func HexStringToRGB(hexstr string) (int, int, int, error) {
	c, err := hex.DecodeString(hexstr) // expected something like: FF0000

	if err != nil {
		return 0, 0, 0, err
	}

	if len(c) != 3 {
		return 0, 0, 0, errors.New("invalid hex string")
	}

	return int(c[0]), int(c[1]), int(c[2]), nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	client := twitch.NewClient("tutkuofnight", os.Getenv("PASS"))
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		r, g, b, _ := HexStringToRGB(message.User.Color[1:])
		color.RGB(r, g, b).Printf("%s:", message.User.DisplayName)
		fmt.Printf("%s\n", message.Message)
		if message.Message == "!ping" {
			client.Say(message.Channel, "pong")
		}
	})

	client.Join("tutkuofnight")
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
