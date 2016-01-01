package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/websocket"
)

// RTMInfo represents a root JSON response of GET /api/rtm.start.
type RTMInfo struct {
	URL      string    `json:"url"`
	Channels []Channel `json:"channels"`
	Users    []User    `json:"users"`
}

// Channel is a JSON object that represents a channel.
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// User is a JSON object that represents a user.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Event is a JSON object that represents a event.
type Event struct {
	ID      string `json:"id"`
	UserID  string `json:"user"`
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	Text    string `json:"text"`
	Ts      string `json:"ts"`
}

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
		log.Fatal(err)
	}

	ws, err := websocket.Dial(info.URL, "", "https://slack.com/")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(ws)
}
