package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// RTMInfo represents a root JSON response of GET /api/rtm.start.
type RTMInfo struct {
	RawURL   string    `json:"url"`
	Channels []Channel `json:"channels"`
	Users    []User    `json:"users"`
}

// URL returns url.URL from RTMInfo.RawURL
func (info *RTMInfo) URL() (*url.URL, error) {
	URL, err := url.Parse(info.RawURL)
	if err != nil {
		return nil, err
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

// StartRTM returns RTMInfo
func StartRTM() (*RTMInfo, error) {
	req, err := http.NewRequest("GET", "https://slack.com/api/rtm.start", nil)
	if err != nil {
		return nil, err
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, errors.New("TOKEN is not found in environment varables")
	}

	query := url.Values{}
	query.Add("token", token)
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info RTMInfo
	err = json.Unmarshal(bytes, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func main() {
	info, err := StartRTM()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(info)
}
