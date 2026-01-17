package server

import (
	"fmt"
	"go-redis/internal/store"
	"net"
)

var redisStore = store.New()

func Start(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	fmt.Println("MiniRedis listening on", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go handleConn(conn)
	}
}
