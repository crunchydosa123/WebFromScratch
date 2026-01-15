package server

import (
	"bufio"
	"fmt"
	"go-redis/internal/resp"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		cmd, err := resp.ParseArray(reader)
		if err != nil {
			return
		}

		fmt.Println("Command:", cmd)

		writer.WriteString("+OK\r\n")
		writer.Flush()
	}
}
