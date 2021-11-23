package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"example.com/main/util"
)

// there's definitely a better way but this works for now
var callBackUrl string = os.Getenv("CALLBACKURL") + "/twitch"

type Slack struct {
	l *log.Logger
}

func NewSlack(l *log.Logger) *Slack {
	return &Slack{l}
}

// Handler interface
func (s *Slack) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	s.l.Println("Slack endpoint hit")

	if r.Method == http.MethodPost {

		// handles slash command
		err := r.ParseForm()
		if err != nil {
			s.l.Println("Error reading form data")
		}

		a := r.FormValue("text")

		command := r.FormValue("command")

		if command != "" {

			// gets channel id
			cID := util.GetUserInfo(a, "id")

			s.l.Println("id-name debug: ", cID)

			if command == "/twitch-add" {

				if cID != "" {

					// call twitch method
					s.l.Println("Sending event subscription request")

					err := util.SendSubRequest(cID, callBackUrl)

					if err != nil {
						fmt.Fprintf(rw, "Error adding event notification")
					} else {

						rw.Header().Set("Content-Type", "application/json")

						resp := &util.SlashResponse{
							ResponseType: "in_channel",
							Text:         "Added event notifications for " + a,
						}

						resp.RespToJSON(rw)
					}

				} else {
					fmt.Fprintf(rw, "User %s not found", a)
				}
			}

			if command == "/twitch-delete" {

				if cID != "" {

					// list active subs
					ActiveSubs := util.GetActiveSubs(s.l)

					var msg string

					// loop through and delete any that match channel name
					for i := 0; i < len(ActiveSubs.Data); i++ {
						if ActiveSubs.Data[i].Condition.BroadcasterUserID == cID {

							// delete this subscription
							err := util.DeleteSub(s.l, ActiveSubs.Data[i].ID)
							if err != nil {
								fmt.Fprintf(rw, "Error deleting sub: %s ", err)
							} else {
								msg = "Deleted event notifications for " + a
							}
						}
					}

					// poorly checking if a msg is null
					if msg == "" {
						msg = "No active subscriptions found for user " + a
					}

					rw.Header().Set("Content-Type", "application/json")

					resp := &util.SlashResponse{
						ResponseType: "in_channel",
						Text:         msg,
					}

					resp.RespToJSON(rw)
				}
			}

			if command == "/twitch-list" {

				// lists active subs by name
				subs := util.GetActiveSubs(s.l)

				mData := subs.GetActiveSubNames(s.l)
				mDataS := strings.Join(mData, " ")

				rw.Header().Set("Content-Type", "application/json")

				resp := &util.SlashResponse{
					ResponseType: "in_channel",
					Text:         "Active event notifications: " + mDataS,
				}

				resp.RespToJSON(rw)

			}
		}
	}
}
