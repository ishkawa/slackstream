package main

import "net/url"

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
	ID        string `json:"id"`
	ChannelID string `json:"channel"`
	UserID    string `json:"user"`
	Type      string `json:"type"`
	Subtype   string `json:"subtype"`
	Text      string `json:"text"`
	Ts        string `json:"ts"`
}

// Team is a JSON object that represents a team.
type Team struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

// RTMInfo represents a root JSON response of GET /api/rtm.start.
type RTMInfo struct {
	RawURL   string    `json:"url"`
	Team     Team      `json:"team"`
	Channels []Channel `json:"channels"`
	Users    []User    `json:"users"`
}

// URL returns url.URL from RTMInfo.RawURL.
func (info *RTMInfo) URL() (*url.URL, error) {
	URL, err := url.Parse(info.RawURL)
	if err != nil {
		return nil, err
	}

	return URL, nil
}
