package pubsub

import (
	"net"
	"testing"
)

func readAll(t *testing.T, conn net.Conn) string {
	t.Helper()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("read failed %w", err)
	}

	return string(buf[:n])
}
