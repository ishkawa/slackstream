package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/websocket"
)

// RTMInfo represents a root JSON response of GET /api/rtm.start.
type RTMInfo struct {
	RawURL   string    `json:"url"`
	Channels []Channel `json:"channels"`
	Users    []User    `json:"users"`
}

// URL returns url.URL from RTMInfo.RawURL.
func (info *RTMInfo) URL() (*url.URL, error) {
	URL, err := url.Parse(info.RawURL)
	if err != nil {
		return nil, err
	}

	comps := strings.Split(URL.Host, ":")
	if err != nil {
		return nil, err
	}

	// Append port if it is missing because net/websocket needs port for wss.
	if len(comps) == 1 {
		URL.Host += ":443"
	}

	return URL, nil
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

	log.Println(ws)
}
