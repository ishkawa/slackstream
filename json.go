package main

import (
	"errors"
	"fmt"
	"net/url"
)

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

// Message represents a message in RTM.
type Message struct {
	Team    *Team
	Event   *Event
	User    *User
	Channel *Channel
}

// NewMessage creates new instance of Message.
func NewMessage(info *RTMInfo, event *Event) (*Message, error) {
	msg := &Message{Team: &info.Team, Event: event}

	for _, user := range info.Users {
		if user.ID == event.UserID {
			msg.User = &user
		}
	}

	for _, channel := range info.Channels {
		if channel.ID == event.ChannelID {
			msg.Channel = &channel
		}
	}

	if msg.User == nil {
		return nil, errors.New("Unknown user " + event.UserID)
	}

	if msg.Channel == nil {
		return nil, errors.New("Unknown channel " + event.ChannelID)
	}

	return msg, nil
}

// Text returns formatted text to display.
func (msg Message) Text() string {
	return fmt.Sprintf("[%s:%s:%s] %s", msg.Team.Domain, msg.Channel.Name, msg.Channel.Name, msg.Event.Text)
}
