package pubsub

import (
	"bufio"
	"net"
	"strconv"
	"sync"
)

type PubSub struct {
	mu       sync.RWMutex
	channels map[string]map[net.Conn]bool
}

func NewPubSub() *PubSub {
	return &PubSub{
		channels: make(map[string]map[net.Conn]bool),
	}
}

func (ps *PubSub) Subscribe(channel string, conn net.Conn) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, ok := ps.channels[channel]; !ok {
		ps.channels[channel] = make(map[net.Conn]bool)
	}
	ps.channels[channel][conn] = true
}

func (ps *PubSub) UnsubscribeAll(conn net.Conn) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, subs := range ps.channels {
		delete(subs, conn)
	}
}

func (ps *PubSub) Publish(channel, message string) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	subs, ok := ps.channels[channel]
	if !ok {
		return 0
	}

	for conn := range subs {
		writePubSubMessage(conn, channel, message)
	}

	return len(subs)
}

func (ps *PubSub) Unsubscribe(channel string, conn net.Conn) int {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	_, ok := ps.channels[channel]
	if !ok {
		return 0
	}

	ps.channels[channel][conn] = false
	return 1
}

func writePubSubMessage(conn net.Conn, channel, message string) {
	writer := bufio.NewWriter(conn)

	writer.WriteString("*3\r\n")
	writer.WriteString("$7\r\nmessage\r\n")
	writer.WriteString("$" + strconv.Itoa(len(channel)) + "\r\n")
	writer.WriteString(channel + "\r\n")
	writer.WriteString("$" + strconv.Itoa(len(message)) + "\r\n")
	writer.WriteString(message + "\r\n")
	writer.Flush()
}
