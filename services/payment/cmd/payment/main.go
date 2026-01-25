package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	nc.Subscribe("orders.process_payment", func(msg *nats.Msg) {
		log.Println("Processing payment for order:", string(msg.Data))

		nc.Publish("payment.succeeded", msg.Data)
	})

	select {}
}
