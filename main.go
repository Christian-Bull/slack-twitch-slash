package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/main/handlers"
	"example.com/main/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	l := log.New(os.Stdout, "slack-api", log.LstdFlags)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	bindAddr := fmt.Sprintf(":%s", port)

	// get twitch auth
	if os.Getenv("BEARERTOKEN") == "" {
		twichApiKey := util.RequestToken(l, os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))

		os.Setenv("BEARERTOKEN", twichApiKey)
	}

	// post test connection message
	err := util.PostMessage(
		l,
		util.CreateMessage("Connected :ratjam:", "C9M568FA4"),
	)
	if err != "" {
		l.Fatal("Error posting slack message: ", err)
	}

	// Create and register handlers
	sh := handlers.NewSlack(l)
	th := handlers.NewTwitch(l)
	hh := handlers.NewHello(l)

	// prometheus metrics
	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/slack", sh)
	http.Handle("/twitch", th)
	http.Handle("/rat", hh)
	http.Handle("/", hh)

	l.Printf("Starting server on port %s", port)
	l.Fatal(http.ListenAndServe(bindAddr, nil))
}
