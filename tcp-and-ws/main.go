package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const (
	network = "tcp"
	port    = ":8888"
)

func main() {
	fmt.Println("listening ", port)
	listener, err := net.Listen(network, port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("connect failed")
			fmt.Println(err)
			return
		}

		go func(conn net.Conn) {
			fmt.Println(conn)
			buf := make([]byte, 1024)
			for {
				n, err := conn.Read(buf)
				if err != nil {
					if err == io.EOF {
						fmt.Println("client disconnected")
						break
					}
					fmt.Println(err)
				}
				fmt.Println(string(buf[:n]))
			}
		}(conn)

	}

}
