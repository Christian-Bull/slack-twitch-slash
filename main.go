package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/main/handlers"
)

func main() {

	l := log.New(os.Stdout, "slack-api", log.LstdFlags)

	port := os.Getenv("PORT")
	if port == "" {
		port = "443"
	}

	bindAddr := fmt.Sprintf(":%s", port)

	// Create and register handlers
	sh := handlers.NewSlack(l)
	th := handlers.NewTwitch(l)
	hh := handlers.NewHello(l)

	http.Handle("/slack", sh)
	http.Handle("/twitch", th)
	http.Handle("/rat", hh)

	l.Printf("Starting server on port %s", port)
	l.Fatal(http.ListenAndServeTLS(bindAddr, "creds/localhost.crt", "creds/localhost.key", nil))
}
