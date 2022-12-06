package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var twitchAuthUrl string = "https://id.twitch.tv/oauth2/token"

type Auth struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func RequestToken(l *log.Logger, clientID string, clientSecret string) string {
	client := &http.Client{}
	auth := &Auth{}

	req, err := http.NewRequest("POST", twitchAuthUrl, nil)
	if err != nil {
		fmt.Println(err)
	}

	q := req.URL.Query()
	q.Add("client_id", clientID)
	q.Add("client_secret", clientSecret)
	q.Add("grant_type", "client_credentials")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getting auth")
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&auth)
	if err != nil {
		l.Println("Error decoding stream info ", err)
	}

	return auth.AccessToken
}
