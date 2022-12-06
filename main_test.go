package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"example.com/main/util"
)

func TestSlackConnection(t *testing.T) {
	l := log.New(os.Stdout, "slack-api-test", log.LstdFlags)

	tstMsg := util.CreateMessage("Test message :ratjam:", "C9M568FA4")

	sendMsg := util.PostMessage(l, tstMsg)
	if sendMsg != "" {
		t.Errorf("Slack Connection error: %s", sendMsg)
	}
}

func TestTwitchResponse(t *testing.T) {
	l := log.New(os.Stdout, "slack-api-test", log.LstdFlags)

	// get bearer token
	if os.Getenv("BEARERTOKEN") == "" {
		twichApiKey := util.RequestToken(l, os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))

		os.Setenv("BEARERTOKEN", twichApiKey)
	}

	// gets stream info for example user twitch
	s := util.GetStreamInfo(l, "136355877")

	var msgText string

	if len(s.Data) > 0 {
		msgText = fmt.Sprintf("%s is now live! Game: %s\n`%s`\n%s",
			"test",
			s.Data[0].GameName,
			s.Data[0].Title,
			"https://www.twitch.tv/"+s.Data[0].UserName)
	} else {
		msgText = fmt.Sprintf("%s is now live! \n %s",
			"Faux",
			"https://www.twitch.tv/Faux")
	}

	t.Logf(msgText)

}
