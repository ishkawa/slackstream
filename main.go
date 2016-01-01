package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

// GetRTMInfo returns RTMInfo
func GetRTMInfo() (*RTMInfo, error) {
	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, errors.New("TOKEN is not found in environment varables")
	}

	resp, err := http.PostForm("https://slack.com/api/rtm.start", url.Values{"token": {token}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info RTMInfo
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func main() {
	info, err := GetRTMInfo()
	if err != nil {
		log.Fatalln("Could not get RTM info:", err)
	}

	URL, err := info.URL()
	if err != nil {
		log.Fatalln("Invalid WebSocket URL:", err)
	}

	ws, err := websocket.Dial(URL.String(), "", "https://slack.com/")
	if err != nil {
		log.Fatalln("Could not establish WebSocket connection:", err)
	}

	go func() {
		for {
			time.Sleep(time.Minute)
			websocket.JSON.Send(ws, map[string]interface{}{
				"id":   1234,
				"type": "ping",
			})
		}
	}()

	for {
		var event Event
		err := websocket.JSON.Receive(ws, &event)
		if err != nil {
			log.Println("Failed to parse JSON:", err)
			continue
		}

		userName := "unknown"
		for _, user := range info.Users {
			if user.ID == event.UserID {
				userName = user.Name
			}
		}

		channelName := "unknown"
		for _, channel := range info.Channels {
			if channel.ID == event.ChannelID {
				channelName = channel.Name
			}
		}

		log.Printf("[%s:%s] %s", channelName, userName, event.Text)
	}
}
