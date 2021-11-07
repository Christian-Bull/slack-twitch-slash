package util

import (
	"log"
	"os"

	"github.com/slack-go/slack"
)

// I guess this will just be for posting messages
// Although I should move all the slash command stuff here as well
// since at some point there will be more than one command

// Message is the struct used to format a slack message
type Message struct {
	message string
	channel string
	status  string
}

func CreateMessage(message string, channel string) Message {
	return Message{
		message: message,
		channel: channel,
		status:  "",
	}
}

func PostMessage(l *log.Logger, m Message) error {
	var (
		retries int = 3
		err     error
	)

	api := slack.New("Bearer " + os.Getenv("SLACKAPIKEY"))

	// retry slack post until it hits the retry limit or is successful
	for retries > 0 {
		msgID, _, _, err := api.SendMessage(
			m.channel,
			slack.MsgOptionText(m.message, false),
		)
		if err != nil {
			l.Println("Error posting message: retry:", retries, err)
			retries--
		} else {
			l.Println("Sent message to: ", m.channel, msgID)
			break
		}

	}
	return err
}