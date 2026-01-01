package main

import (
	"eventure/libs/events"
	nss "eventure/libs/nats"
	"eventure/services/payment/internal/ports"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc := nss.ConnectNATS()
	defer nc.Close()

	handler := ports.NewPaymentHandler(nc)

	nc.Subscribe("payment.charge", func(msg *nats.Msg) {
		evt, err := events.UnmarshalPaymentCharge(msg.Data)
		if err != nil {
			log.Println("invalid payment.charge message:", err)
			return
		}
		handler.Charge(evt.OrderID)
	})

	log.Println("Payment service running...")
	select {}
}
