package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/c00/botman/models"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func SetupSlackbot() error {

	//Get the user Id
	response, err := client.AuthTest()
	if err != nil {
		return err
	}
	userId = response.UserID

	// Create a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())
	// Make this cancel called properly in a real program , graceful shutdown etc
	defer cancel()

	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		// Create a for loop that selects either the context cancellation or the events incomming
		for {
			select {
			// inscase context cancel is called exit the goroutine
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-socketClient.Events:
				switch event.Type {
				case socketmode.EventTypeEventsAPI:
					eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					socketClient.Ack(*event.Request)
					err := handleEvent(eventsAPIEvent)
					if err != nil {
						fmt.Println("Error handling Slack Event:", err)
					}
				}
			}
		}
	}(ctx, client, socketClient)

	return socketClient.Run()
}

func handleEvent(event slackevents.EventsAPIEvent) error {
	switch event.Type {
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent

		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			if ev.User == userId {
				return nil
			}

			//Direct message?
			if ev.ChannelType == "im" {
				fmt.Printf("Handling IM Message: '%v'\n", getSubstring(ev.Text))
				return handleMessage(ev)
			}

			//Check for tag
			if strings.Contains(ev.Text, fmt.Sprintf("<@%v>", userId)) {
				fmt.Printf("Handling Tagged Message: '%v'\n", getSubstring(ev.Text))
				return handleMessage(ev)
			}

			return nil

		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

func handleMessage(ev *slackevents.MessageEvent) error {
	messages := []models.ChatMessage{
		{Role: models.ChatMessageRoleSystem, Content: botmanConf.SystemPrompt},
	}

	ts := ev.TimeStamp
	if ev.ThreadTimeStamp != "" {
		ts = ev.ThreadTimeStamp
		//Get entire thread
		thread, err := getThread(ev)
		if err != nil {
			return err
		}

		//Add it to messages
		messages = append(messages, thread...)
	} else {
		messages = append(messages, models.ChatMessage{Role: models.ChatMessageRoleUser, Content: cleanMessage(ev.Text)})
	}

	err := client.AddReaction("thought_balloon", slack.NewRefToMessage(ev.Channel, ev.TimeStamp))
	if err != nil {
		fmt.Println("Error adding reaction", err)
		return err
	}

	response := respond(messages)
	client.SendMessage(ev.Channel, slack.MsgOptionText(response, false), slack.MsgOptionTS(ts))

	client.AddReaction("ballot_box_with_check", slack.NewRefToMessage(ev.Channel, ev.TimeStamp))
	client.RemoveReaction("thought_balloon", slack.NewRefToMessage(ev.Channel, ev.TimeStamp))

	return nil
}

func respond(messages []models.ChatMessage) string {
	channel := getBlackHole()
	result := chatter.GetResponse(messages, channel)

	return result
}

func getThread(ev *slackevents.MessageEvent) ([]models.ChatMessage, error) {
	messages, _, _, err := client.GetConversationReplies(&slack.GetConversationRepliesParameters{ChannelID: ev.Channel, Timestamp: ev.ThreadTimeStamp})
	if err != nil {
		return nil, err
	}

	result := make([]models.ChatMessage, 0, len(messages))

	current := ""

	for _, msg := range messages {
		if msg.User == userId {
			if current != "" {
				result = append(result, models.ChatMessage{Role: models.ChatMessageRoleUser, Content: current})
				current = ""
			}
			result = append(result, models.ChatMessage{Role: models.ChatMessageRoleAssistant, Content: msg.Text})
		} else {
			addition := cleanMessage(msg.Text)
			if addition != "" {
				current += fmt.Sprintf("<@%v> says:\n%v\n\n", msg.User, addition)
			}
		}
	}

	fmt.Printf("Got thread with %v messages\n", len(messages))

	//Add New message
	addition := cleanMessage(ev.Text)
	if addition != "" {
		current += fmt.Sprintf("<@%v> says:\n%v\n\n", ev.User, addition)
	}

	if current != "" {
		result = append(result, models.ChatMessage{Role: models.ChatMessageRoleUser, Content: current})
	}

	return result, nil
}

func cleanMessage(message string) string {
	message = strings.ReplaceAll(message, fmt.Sprintf("<@%v>", userId), "")
	message = strings.TrimSpace(message)
	return message
}

func getSubstring(msg string) string {
	if len(msg) == 0 {
		return ""
	}

	if len(msg) > 10 {
		return fmt.Sprintf("%v...", msg[:7])
	}

	return msg
}
