package server

import (
	"fmt"
	"go-redis/internal/aof"
	pubsub "go-redis/internal/pub-sub"
	"go-redis/internal/store"
	"net"
)

var redisStore = store.New()
var pubSub = pubsub.NewPubSub()
var aofLog *aof.AOF
var aofEnabled = false

func Start(addr string) {
	/*aofLog, err := aof.New("appendonly.aof")
	if err != nil {
		panic(err)
	}

	aofEnabled = false
	err = aofLog.Replay(func(cmd []string) {
		executeCommand(cmd)
	})
	aofEnabled = true

	if err != nil {
		panic(err)
	}*/

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

/*func executeCommand(cmd []string) {
	modified := applyCommand(cmd)

	if aofEnabled && modified {
		aofLog.Append(cmd)
	}
}*/
