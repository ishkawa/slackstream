package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/websocket"
)

// RTMConn wraps *websocket.Conn
type RTMConn struct {
	Info *RTMInfo
	Ws   *websocket.Conn
}

// NewRTMConn returns new RTMConn.
// This method fetches RTM info and open new WebSocket connection.
// https://api.slack.com/methods/rtm.start
func NewRTMConn(token string) (*RTMConn, error) {
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

	URL, err := info.URL()
	if err != nil {
		return nil, err
	}

	ws, err := websocket.Dial(URL.String(), "", "https://slack.com/")
	if err != nil {
		return nil, err
	}

	conn := &RTMConn{
		Info: &info,
		Ws:   ws,
	}

	return conn, nil
}

// Run starts handling events and queue messages into the passed channel.
func (conn *RTMConn) Run(msgs chan *Message) {
	id := 1
	timer := make(chan bool)
	go startTimer(timer)

	events := make(chan Event)
	go pipeEvent(conn.Ws, events)

	for {
		select {
		case <-timer:
			// id must be unique for a connection.
			id++
			conn.Ping(id)

		case event := <-events:
			msg, err := NewMessage(conn.Info, &event)
			if err != nil {
				// TODO: Log erros
				continue
			}

			msgs <- msg
		}
	}
}

// Ping sends ping signal to websocket.Conn.
func (conn *RTMConn) Ping(id int) {
	websocket.JSON.Send(conn.Ws, map[string]interface{}{
		"id":   id,
		"type": "ping",
	})
}

func pipeEvent(ws *websocket.Conn, events chan Event) {
	for {
		var event Event
		err := websocket.JSON.Receive(ws, &event)
		if err != nil {
			// TODO: Handle error
			continue
		}

		events <- event
	}
}

func startTimer(timer chan bool) {
	for {
		time.Sleep(time.Minute)
		timer <- true
	}
}
