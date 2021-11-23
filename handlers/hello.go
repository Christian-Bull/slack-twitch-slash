package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (s *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.l.Println("home endpoint hit")

	fmt.Fprintf(rw, "RAT GANG X PBR")
}
