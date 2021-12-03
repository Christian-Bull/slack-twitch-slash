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

			// return response first
			rw.WriteHeader(http.StatusOK)

			// double check it's the correct one
			n := &util.EventNotification{}

			nDecoder := json.NewDecoder(r.Body)
			err := nDecoder.Decode(&n)
			if err != nil {
				fmt.Println("Error decoding body", err)
			}

			if n.Subscription.Type == "stream.online" {

				var msgText string

				// get stream infos
				sInfo := util.GetStreamInfo(s.l, n.Event.UserID)

				// check if stream info was returned
				if len(sInfo.Data) > 0 {
					msgText = fmt.Sprintf("%s is now live! Game: %s\n`%s`\n%s",
						n.Event.BroadcasterUserName,
						sInfo.Data[0].GameName,
						sInfo.Data[0].Title,
						"https://www.twitch.tv/"+sInfo.Data[0].UserName)
				} else {
					msgText = fmt.Sprintf("%s is now live! \n %s",
						n.Event.BroadcasterUserName,
						"https://www.twitch.tv/"+n.Event.BroadcasterUserName)
				}

				// post notification to slack
				fmt.Println(n.Event.BroadcasterUserName)
				msg := util.CreateMessage(
					msgText,
					os.Getenv("CHANNEL"),
				)

				err := util.PostMessage(s.l, msg)
				if err != "" {
					s.l.Println("Error posting slack message: ", err)
				}

			} else {
				fmt.Println("Event type: ", n.Subscription.Type, n.Event.BroadcasterUserName)
			}
		}

	}
}
