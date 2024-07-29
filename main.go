package main

import (
	"fmt"
	"os"

	"github.com/c00/botman/botman"
	"github.com/c00/botman/config"
	"github.com/c00/botman/models"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

var chatter models.Chatter
var botmanConf config.BotmanConfig
var userId string
var client *slack.Client
var socketClient *socketmode.Client
var verbose bool

func main() {
	godotenv.Load(".env")

	client = slack.New(os.Getenv("SLACK_BOT_TOKEN"), slack.OptionAppLevelToken(os.Getenv("SLACK_APP_TOKEN")))
	socketClient = socketmode.New(client)

	verbose = boolFromEnv("BOTMAN_VERBOSE", false)

	fmt.Println("Starting Botman Slackbot...")

	botmanConf = config.LoadFromEnv()
	chatter = botman.GetChatter(botmanConf)

	err := SetupSlackbot()
	if err != nil {
		fmt.Println("Error settig up Slackbot:", err)
	}
}

func boolFromEnv(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	if val == "1" || val == "true" {
		return true
	}

	return false
}

// Returns a channel that ignores all values sent to it.
// Used to ignore the streaming portion of the response.
func getBlackHole() chan string {
	ch := make(chan string)

	go func() {
		for range ch {
			//Just do nothing.
		}
	}()

	return ch
}
