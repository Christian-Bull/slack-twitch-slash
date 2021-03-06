package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

	payLoadBuf := new(bytes.Buffer)
	json.NewEncoder(payLoadBuf).Encode(body)
	req, err := http.NewRequest(http.MethodPost, postEndpoint, payLoadBuf)

	req.Header.Add("Client-ID", os.Getenv("CLIENT_ID"))
	req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARERTOKEN"))
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

// provide either a channel name or id, it will return the opposite
func GetUserInfo(input string, outType string) string {
	var url string
	var out string

	if outType == "id" {
		// provide channel name
		url = "https://api.twitch.tv/helix/users?login=" + input
	} else if outType == "name" {
		url = "https://api.twitch.tv/helix/users?id=" + input
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("Client-ID", os.Getenv("CLIENT_ID"))
	req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARERTOKEN"))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error: ", err)
	}

	if res.StatusCode == 200 {

		defer res.Body.Close()

		userInfo := &UserInfo{}
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&userInfo)
		if err != nil {
			fmt.Println(err)
		}

		if outType == "id" {
			out = userInfo.Data[0].ID
		} else if outType == "name" {
			out = userInfo.Data[0].DisplayName
		}
	}

	return out
}

type ActiveSubs struct {
	Total int `json:"total"`
	Data  []struct {
		ID        string `json:"id"`
		Status    string `json:"status"`
		Type      string `json:"type"`
		Version   string `json:"version"`
		Condition struct {
			BroadcasterUserID string `json:"broadcaster_user_id"`
		} `json:"condition"`
		CreatedAt time.Time `json:"created_at"`
		Transport struct {
			Method   string `json:"method"`
			Callback string `json:"callback"`
		} `json:"transport"`
		Cost int `json:"cost"`
	} `json:"data"`
	MaxTotalCost int `json:"max_total_cost"`
	TotalCost    int `json:"total_cost"`
	Pagination   struct {
	} `json:"pagination"`
}

func GetActiveSubs(l *log.Logger) *ActiveSubs {

	subs := &ActiveSubs{}

	url := "https://api.twitch.tv/helix/eventsub/subscriptions"
	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("Client-ID", os.Getenv("CLIENT_ID"))
	req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARERTOKEN"))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error: ", err)
	}

	if res.StatusCode == 200 {
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&subs)
		if err != nil {
			l.Println("Error decoding subs: ", err)
		}

	} else {
		defer res.Body.Close()

		respBody, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response", err)
		}

		fmt.Println(string(respBody))

		l.Println(res.StatusCode)
	}

	return subs
}

func DeleteSub(l *log.Logger, SubID string) error {
	var err error

	// send a sub delete request
	url := "https://api.twitch.tv/helix/eventsub/subscriptions?id=" + SubID
	req, err := http.NewRequest(http.MethodDelete, url, nil)

	req.Header.Add("Client-ID", os.Getenv("CLIENT_ID"))
	req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARERTOKEN"))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error: ", err)
	}

	if res.StatusCode != 204 {
		err = errors.New("Error deleting subscription")
	} else if res.StatusCode == 404 {
		err = errors.New("Subscription for channel not found")
	}

	return err
}

func (a *ActiveSubs) subNamesToList() []string {
	var aList []string

	for i := 0; i < len(a.Data); i++ {
		aList = append(aList, a.Data[i].Condition.BroadcasterUserID)
	}
	return aList
}

func channelIDsToName(listOfIDs []string) []string {
	var out []string

	for i := 0; i < len(listOfIDs); i++ {
		a := GetUserInfo(listOfIDs[i], "name")
		out = append(out, a)
	}

	return out
}

func (a *ActiveSubs) GetActiveSubNames(l *log.Logger) []string {
	s := a.subNamesToList()

	return channelIDsToName(s)
}

type StreamInfo struct {
	Data []struct {
		ID           string    `json:"id"`
		UserID       string    `json:"user_id"`
		UserLogin    string    `json:"user_login"`
		UserName     string    `json:"user_name"`
		GameID       string    `json:"game_id"`
		GameName     string    `json:"game_name"`
		Type         string    `json:"type"`
		Title        string    `json:"title"`
		ViewerCount  int       `json:"viewer_count"`
		StartedAt    time.Time `json:"started_at"`
		Language     string    `json:"language"`
		ThumbnailURL string    `json:"thumbnail_url"`
		TagIds       []string  `json:"tag_ids"`
		IsMature     bool      `json:"is_mature"`
	} `json:"data"`
	Pagination struct {
	} `json:"pagination"`
}

func GetStreamInfo(l *log.Logger, userID string) *StreamInfo {

	sInfo := &StreamInfo{}

	url := "https://api.twitch.tv/helix/streams?user_id=" + userID
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("Client-ID", os.Getenv("CLIENT_ID"))
	req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARERTOKEN"))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		l.Println("Stream info request error: ", err)
	}

	if res.StatusCode == 200 {
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&sInfo)
		if err != nil {
			l.Println("Error decoding stream info ", err)
		}

	} else {
		l.Println("Twitch stream info error: ", res.StatusCode)
	}

	return sInfo

}

// cleans up subscriptions
// resubscribes to any subscription that's gone inactive
// due to notification_failures_exceeded
func (a *ActiveSubs) Resubscribe(l *log.Logger) {

	// loop through active subs
	for i := 0; i < len(a.Data); i++ {

		// check the status
		if a.Data[i].Status == "notification_failures_exceeded" {

			l.Printf("Resubscribe debug: %s %s %s %s",
				a.Data[i].Status,
				a.Data[i].ID,
				a.Data[i].Condition.BroadcasterUserID,
				a.Data[i].Transport.Callback,
			)

			l.Printf("%+v\n", a.Data[i])

			// grab the broadcaster id
			broadcasterID := a.Data[i].Condition.BroadcasterUserID

			// delete the old sub
			err := DeleteSub(l, a.Data[i].ID)
			if err != nil {
				l.Println("Error deleting subscription: " + a.Data[i].ID)
			}

			// resubscribe
			channelName := GetUserInfo(broadcasterID, "name")

			// it's possible a channel has been deleted
			if channelName != "" {

				err = SendSubRequest(channelName, os.Getenv("CALLBACKURL")+"/twitch")
				if err != nil {
					l.Println("Error subbing to: " + channelName + " resubscribe debug")
				}
			}
		}
	}
}
