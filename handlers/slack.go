package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
		if a == "" {
			s.l.Println("No text form key found")
		}

		// gets channel name
		cName := util.GetUserID(a)

		if cName != "" {

			// call twitch method
			s.l.Println("Sending event subscription request")
			s.l.Println(callBackUrl)
			util.SendSubRequest(cName, callBackUrl)

			fmt.Fprintf(rw, "Subscribed to %s for twitch notifications", a)
		}

	}
}
