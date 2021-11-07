package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/main/util"
)

type Twitch struct {
	l *log.Logger
}

func NewTwitch(l *log.Logger) *Twitch {
	return &Twitch{l}
}

// Handler interface
func (s *Twitch) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	s.l.Println("Twitch endpoint hit")

	if r.Method == http.MethodPost {

		// check for header if it's a callback verification
		hc := r.Header.Get("Twitch-Eventsub-Message-Type")
		if hc != "" && hc == "webhook_callback_verification" {

			// confirm challenge
			c := &util.CallbackVerify{}

			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&c)
			if err != nil {
				fmt.Println("Error decoding rbody", err)
			}

			answer := c.Challenge

			fmt.Fprintln(rw, answer)

		}

		// check for a notification header
		nh := r.Header.Get("Twitch-Eventsub-Message-Type")
		if nh != "" && nh == "notification" {

			// double check it's the correct one
			n := &util.EventNotification{}

			nDecoder := json.NewDecoder(r.Body)
			err := nDecoder.Decode(&n)
			if err != nil {
				fmt.Println("Error decoding body", err)
			}

			if n.Subscription.Type == "stream.online" {

				// post notification to slack
				fmt.Println(n.Event.BroadcasterUserName)
				msg := util.CreateMessage(
					n.Event.BroadcasterUserName+" is now live!",
					os.Getenv("CHANNEL"),
				)

				err := util.PostMessage(s.l, msg)
				if err != nil {
					s.l.Println("Error posting slack message: ", err)
				}

			} else {
				fmt.Println("Event type: ", n.Subscription.Type, n.Event.BroadcasterUserName)
			}
		}

	}
}
