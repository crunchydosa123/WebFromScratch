package server

import (
	"bufio"
	"go-redis/internal/resp"
	"net"
	"strconv"
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

			var ttl *int = nil
			if len(cmd) == 5 {
				if strings.ToUpper(cmd[3]) != "EX" {
					resp.WriteError(writer, "syntax error")
					break
				}

				seconds, err := strconv.Atoi(cmd[4])
				if err != nil {
					resp.WriteError(writer, "invalid expire time")
					break
				}

				ttl = &seconds
			}

			redisStore.Set(key, value, ttl)
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
