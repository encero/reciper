package tests

import (
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
)

func RunNatsServer(t *testing.T) (*server.Server, string) {
	s := natsserver.RunRandClientPortServer()

	info := s.PortsInfo(time.Second)

	if len(info.Nats) == 0 {
		t.Fatalf("no nats ports")
	}

	return s, info.Nats[0]
}

func RunAndConnectNats(t *testing.T) (*nats.Conn, string, func()) {
	s, url := RunNatsServer(t)

	conn, err := nats.Connect(url)
	if err != nil {
		t.Fatalf("nats connect: %s", err)
	}

	return conn, url, func() {
		conn.Close()
		s.Shutdown()
	}
}
