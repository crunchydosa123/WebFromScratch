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
				resp.WriteSimpleString(writer, "PONG")
			} else {
				resp.WriteSimpleString(writer, cmd[1])
			}
		case "SET":
			if len(cmd) != 3 {
				resp.WriteError(writer, "wrong number of arguments for 'get'")
				break
			}
			key := cmd[1]
			value := cmd[2]

			redisStore.Set(key, value)
			resp.WriteSimpleString(writer, "OK")
		case "GET":
			if len(cmd) != 2 {
				resp.WriteError(writer, "wrong number of arguments for 'set'")
				break
			}
			key := cmd[1]

			val, ok := redisStore.Get(key)
			if !ok {
				resp.WriteBulkString(writer, nil)
			} else {
				resp.WriteBulkString(writer, &val)
			}
		case "COMMAND":
			writer.WriteString("*0\r\n")
		default:
			resp.WriteError(writer, "unknown command '"+cmd[0]+"'")
		}

		writer.Flush()
	}
}
