package main

import (
	"log"
	"os"
	"time"
)

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalln("TOKEN is not found in environment varables")
	}

	conn, err := OpenRTMConn(token)
	if err != nil {
		log.Fatalln("Could not open RTM connection:", err)
	}

	go func() {
		for {
			time.Sleep(time.Minute)
			conn.Ping()
		}
	}()

	for {
		msg, err := conn.ReceiveMsg()
		if err == nil {
			log.Println(msg)
		}
	}
}
