package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

			s.l.Println(
				"Broadcaster: "+n.Event.BroadcasterUserID,
				"Broadcaster: "+n.Event.BroadcasterUserName,
				"User: "+n.Event.UserID,
				"User: "+n.Event.UserName,
			)

			if n.Subscription.Type == "stream.online" {

				var msgText string

				// get stream infos
				sInfo := util.GetStreamInfo(s.l, n.Event.BroadcasterUserID)

				// check if stream info was returned
				if len(sInfo.Data) > 0 {
					msgText = fmt.Sprintf("%s is now live! Game: %s\n`%s`\n%s",
						n.Event.BroadcasterUserName,
						sInfo.Data[0].GameName,
						sInfo.Data[0].Title,
						"https://www.twitch.tv/"+n.Event.BroadcasterUserName)
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

				// temp for logging how long it takes for stream info to show
				go func() {
					var retries int = 5
					var delay int = 20

					for i := 0; i < retries; i++ {
						sTmp := util.GetStreamInfo(s.l, n.Event.BroadcasterUserID)

						// log a bunch of stuff
						if len(sTmp.Data) > 0 {
							s.l.Println(sTmp.Data[0].StartedAt)
							s.l.Println(sTmp.Data[0].GameName)
							s.l.Println(sTmp.Data[0].Title)
							break
						} else {
							s.l.Println("No stream info for: ", n.Event.BroadcasterUserName)
						}

						time.Sleep(time.Duration(delay) * time.Second)
					}
				}()

			} else {
				fmt.Println("Event type: ", n.Subscription.Type, n.Event.BroadcasterUserName)
			}
		}

	}
}
