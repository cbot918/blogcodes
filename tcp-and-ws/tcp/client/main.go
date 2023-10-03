package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	network = "tcp"
	port    = ":8888"
	wait    = 3 * time.Second
	users   = 5
)

func main() {

	fmt.Printf("launch %d users", users)

	for i := 0; i < users; i++ {
		go runClient(i)
	}

	select {}
}

func runClient(number int) {
	conn, err := net.Dial(network, port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		message := fmt.Sprintf("hello from user %d", number)
		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(wait)
	}
}
