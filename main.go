package main

import (
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	conns := []*RTMConn{}
	for _, token := range strings.Split(os.Getenv("TOKEN"), ",") {
		conn, err := OpenRTMConn(token)
		if err != nil {
			log.Fatalln("Could not open RTM connection:", err)
		}

		conns = append(conns, conn)
	}

	go func() {
		for {
			time.Sleep(time.Minute)

			for _, conn := range conns {
				conn.Ping()
			}
		}
	}()

	ch := make(chan string)

	for _, conn := range conns {
		go func(conn *RTMConn) {
			for {
				time.Sleep(time.Minute)
				conn.Ping()
			}
		}(conn)

		go func(conn *RTMConn) {
			for {
				msg, err := conn.ReceiveMsg()
				if err == nil {
					ch <- msg
				}
			}
		}(conn)
	}

	for msg := range ch {
		log.Println(msg)
	}
}
