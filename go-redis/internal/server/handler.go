package server

import (
	"bufio"
	"go-redis/internal/resp"
	"net"
	"strings"
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

		if len(cmd) == 0 {
			resp.WriteError(writer, "empty command")
			writer.Flush()
			continue
		}

		command := strings.ToUpper(cmd[0])

		switch command {
		case "PING":
			if len(cmd) == 1 {
				resp.WriteSimpleString(writer, "PONG") //implement well
			} else {
				resp.WriteSimpleString(writer, cmd[1])
			}
		default:
			resp.WriteError(writer, "unknown command '"+cmd[0]+"'")
		}

		writer.Flush()
	}
}
