package handlers

import (
	"fmt"
	"log"
	"net/http"
)

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

		// call twitch method
		// check if already subscribed
		// subscribe if not

		fmt.Fprintf(rw, "Hi there, I love %s!", a)
	}
}
