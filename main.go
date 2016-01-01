package main

import (
	"log"
	"time"
)

func main() {
	conn, err := OpenRTMConn()
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
