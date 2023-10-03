package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
)

const (
	network = "tcp"
	port    = "localhost:8889"
)

func main() {
	fmt.Println("listeninging: ", port)
	listener, err := net.Listen(network, port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(conn)
		}

		go handleConnection(conn)

	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("client disconnetced")
				return
			}
			fmt.Println("read socket failed")
			return
		}
		data := buf[:n]
		fmt.Println(string(buf[:n]))

		// 升級連線的核心函式
		key := getWebSecKey(data)
		retKey := getReturnSec(key)

		// 把 response 字串透過 conn 寫回去給 client 端, 完成升級連線
		// 此 socket 之後就可以用 websocket 的方式 傳送資料(也就是需要按照ws格式編解碼)
		response := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n\r\n", retKey)
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("checkout web for socket connected")
	}
}

func getWebSecKey(data []byte) string {
	pattern := `Sec-WebSocket-Key: ([^\r\n]+)`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(string(data))
	return strings.TrimSpace(match[1])
}

func getReturnSec(webSecSocketkey string) string {
	var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	h := sha1.New()
	h.Write([]byte(webSecSocketkey))
	h.Write(keyGUID)
	secWebSocketAccept := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return secWebSocketAccept
}
