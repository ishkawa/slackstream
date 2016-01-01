package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/websocket"
)

// RTMConn wraps *websocket.Conn
type RTMConn struct {
	Info *RTMInfo
	Ws   *websocket.Conn
}

// Ping sends ping signal to websocket.Conn.
func (conn *RTMConn) Ping() {
	websocket.JSON.Send(conn.Ws, map[string]interface{}{
		"id":   1234,
		"type": "ping",
	})
}

// ReceiveMsg waits WebSocket chunk and returns formatted message.
func (conn *RTMConn) ReceiveMsg() (string, error) {
	var event Event
	err := websocket.JSON.Receive(conn.Ws, &event)
	if err != nil {
		return "", err
	}

	userName := ""
	for _, user := range conn.Info.Users {
		if user.ID == event.UserID {
			userName = user.Name
		}
	}

	channelName := ""
	for _, channel := range conn.Info.Channels {
		if channel.ID == event.ChannelID {
			channelName = channel.Name
		}
	}

	if userName == "" {
		return "", errors.New("Unknown user " + event.UserID)
	}

	if channelName == "" {
		return "", errors.New("Unknown channel " + event.ChannelID)
	}

	message := fmt.Sprintf("[%s:%s] %s", channelName, userName, event.Text)

	return message, nil
}

// OpenRTMConn returns new RTMConn.
// This method fetches RTM info and open new WebSocket connection.
func OpenRTMConn() (*RTMConn, error) {
	info, err := FetchRTMInfo()
	if err != nil {
		return nil, err
	}

	URL, err := info.URL()
	if err != nil {
		return nil, err
	}

	ws, err := websocket.Dial(URL.String(), "", "https://slack.com/")
	if err != nil {
		return nil, err
	}

	conn := &RTMConn{
		Info: info,
		Ws:   ws,
	}

	return conn, nil
}

// FetchRTMInfo fetches RTM info from Slack server.
// https://api.slack.com/methods/rtm.start
func FetchRTMInfo() (*RTMInfo, error) {
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
