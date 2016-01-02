package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	msgs := make(chan *Message)
	for _, token := range strings.Split(os.Getenv("TOKEN"), ",") {
		conn, err := NewRTMConn(token)
		if err != nil {
			log.Fatalln("Could not open RTM connection:", err)
		}

		go conn.Run(msgs)
	}

	for msg := range msgs {
		log.Println(msg.Text())
	}
}
