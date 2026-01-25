package messaging

import "github.com/nats-io/nats.go"

func Connect() (*nats.Conn, error) {
	// NATS connection logic here
	return nats.Connect(nats.DefaultURL)
}
