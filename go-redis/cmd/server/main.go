package main

import "go-redis/internal/server"

func main() {
	server.Start(":6379")
}
