package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	fmt.Println("Connected to server")

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	for {
		fmt.Println("> ")
		text, _ := reader.ReadString('\n')
		conn.Write([]byte(text))

		resp, _ := serverReader.ReadString('\n')
		fmt.Println("server:", resp)
	}
}
