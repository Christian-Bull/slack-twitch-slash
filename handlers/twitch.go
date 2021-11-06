package handlers

import (
	"fmt"
	"log"
	"net/http"
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

	fmt.Fprintf(rw, "Hi there, I love %s!", r.URL)

}
