package main

import (
	"log"
	"net/http"
	"os"

	"example.com/main/handlers"
)

func main() {

	l := log.New(os.Stdout, "slack-api", log.LstdFlags)

	sh := handlers.NewSlack(l)

	http.Handle("/", sh)
	log.Fatal(http.ListenAndServe(":5050", nil))
}
