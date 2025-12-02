package nats

import (
	"log"

	"github.com/nats-io/nats.go"
)

func ConnectNATS() *nats.Conn {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("failed to connect to NATs: %v", err)
	}

	return nc
}
