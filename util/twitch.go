package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Condition struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type Transport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

type SubReq struct {
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	Condition Condition `json:"condition"`
	Transport Transport `json:"transport"`
}

type CallbackVerify struct {
	Subscription struct {
		ID        string    `json:"id"`
		Status    string    `json:"status"`
		Type      string    `json:"type"`
		Version   string    `json:"version"`
		Condition Condition `json:"condition"`
		Transport Transport `json:"transport"`
		CreatedAt time.Time `json:"created_at"`
		Cost      int       `json:"cost"`
	} `json:"subscription"`
	Challenge string `json:"challenge"`
}

type EventNotification struct {
	Subscription struct {
		ID        string    `json:"id"`
		Status    string    `json:"status"`
		Type      string    `json:"type"`
		Version   string    `json:"version"`
		Cost      int       `json:"cost"`
		Condition Condition `json:"condition"`
		Transport Transport `json:"transport"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"subscription"`
	Event struct {
		UserID               string `json:"user_id"`
		UserLogin            string `json:"user_login"`
		UserName             string `json:"user_name"`
		BroadcasterUserID    string `json:"broadcaster_user_id"`
		BroadcasterUserLogin string `json:"broadcaster_user_login"`
		BroadcasterUserName  string `json:"broadcaster_user_name"`
	} `json:"event"`
}

type UserInfo struct {
	Data []struct {
		ID              string    `json:"id"`
		Login           string    `json:"login"`
		DisplayName     string    `json:"display_name"`
		Type            string    `json:"type"`
		BroadcasterType string    `json:"broadcaster_type"`
		Description     string    `json:"description"`
		ProfileImageURL string    `json:"profile_image_url"`
		OfflineImageURL string    `json:"offline_image_url"`
		ViewCount       int       `json:"view_count"`
		CreatedAt       time.Time `json:"created_at"`
	} `json:"data"`
}

var postEndpoint string = "https://api.twitch.tv/helix/eventsub/subscriptions"

func SendSubRequest(channelName string, callbackUrl string) error {

	body := &SubReq{
		Type:    "stream.online",
		Version: "1",
		Condition: Condition{
			BroadcasterUserID: channelName,
		},
		Transport: Transport{
			Method:   "webhook",
			Callback: callbackUrl,
			Secret:   os.Getenv("CLIENT_SECRET"),
		},
	}

	fmt.Println(body)

	payLoadBuf := new(bytes.Buffer)
	json.NewEncoder(payLoadBuf).Encode(body)
	req, err := http.NewRequest(http.MethodPost, postEndpoint, payLoadBuf)

	req.Header.Add("Client-ID", os.Getenv("CLIENT_ID"))
	req.Header.Add("Authorization", os.Getenv("BEARERTOKEN"))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error: ", err)
	}

	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response", err)
	}

	fmt.Println(string(respBody))

	return err
}

func GetUserID(channelName string) string {
	url := "https://api.twitch.tv/helix/users?login=" + channelName
	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("Client-ID", os.Getenv("CLIENT_ID"))
	req.Header.Add("Authorization", os.Getenv("BEARERTOKEN"))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error: ", err)
	}

	defer res.Body.Close()

	userInfo := &UserInfo{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&userInfo)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(userInfo)
	// hopefully there's only one user
	return userInfo.Data[0].ID
}
