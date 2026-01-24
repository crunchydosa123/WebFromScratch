package server

import (
	"bufio"
	"go-redis/internal/resp"
	"net"
	"strconv"
	"strings"
)

func handleConn(conn net.Conn) {
	isSubscriber := false
	defer func() {
		pubSub.UnsubscribeAll(conn)
		conn.Close()
	}()

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

		if isSubscriber &&
			command != "SUBSCRIBE" &&
			command != "UNSUBSCRIBE" &&
			command != "PING" {
			resp.WriteError(writer, "ERR only (un)subscribe allowed in subscriber mode")
			writer.Flush()
			continue
		}

		switch command {
		case "PING":
			if len(cmd) == 1 {
				resp.WriteSimpleString(writer, "PONG")
			} else {
				resp.WriteSimpleString(writer, cmd[1])
			}
		case "SET":
			if len(cmd) != 3 && len(cmd) != 5 {
				resp.WriteError(writer, "wrong number of arguments for 'set'")
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
				resp.WriteError(writer, "wrong number of arguments for 'get'")
				break
			}
			key := cmd[1]

			val, ok := redisStore.Get(key)
			if !ok {
				resp.WriteBulkString(writer, nil)
			} else {
				resp.WriteBulkString(writer, &val)
			}
		case "SUBSCRIBE":
			if len(cmd) != 2 {
				resp.WriteError(writer, "wrong number of arguments for 'subscribe'")
				break
			}

			channel := cmd[1]
			pubSub.Subscribe(channel, conn)
			isSubscriber = true

			writer.WriteString("*3\r\n")
			writer.WriteString("$9\r\nsubscribe\r\n")
			writer.WriteString("$" + strconv.Itoa(len(channel)) + "\r\n" + channel + "\r\n")
			writer.WriteString(":1\r\n")

		case "UNSUBSCRIBE":
			if len(cmd) > 2 {
				resp.WriteError(writer, "wrong number of arguments for 'unsubscribe'")
				break
			}

			if len(cmd) == 2 && cmd[1] != "ALL" {
				channel := cmd[1]
				remaining := pubSub.Unsubscribe(channel, conn)

				writer.WriteString("*3\r\n")
				writer.WriteString("$11\r\nunsubscribe\r\n")
				writer.WriteString("$" + strconv.Itoa(len(channel)) + "\r\n" + channel + "\r\n")
				writer.WriteString(":" + strconv.Itoa(remaining) + "\r\n")
			}

			if len(cmd) == 2 && cmd[1] == "ALL" {
				pubSub.UnsubscribeAll(conn)
				isSubscriber = false
				writer.WriteString("*3\r\n")
				writer.WriteString("$9\r\nunsubscribe all\r\n")
				writer.WriteString(":1\r\n")
			}
		case "PUBLISH":
			if len(cmd) != 3 {
				resp.WriteError(writer, "wrong number of arguments for 'publish'")
				break
			}

			channel := cmd[1]
			message := cmd[2]

			count := pubSub.Publish(channel, message)
			writer.WriteString(":" + strconv.Itoa(count) + "\r\n")
		case "COMMAND":
			writer.WriteString("*0\r\n")
		default:
			resp.WriteError(writer, "unknown command '"+cmd[0]+"'")
		}

		writer.Flush()
	}
}
