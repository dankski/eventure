package main

import (
	"log"

	"eventure/libs/events"
	nss "eventure/libs/nats"
	"eventure/services/inventory/internal/ports"

	"github.com/nats-io/nats.go"
)

func main() {
	nc := nss.ConnectNATS()
	defer nc.Close()

	handler := ports.NewInventoryHandler()

	_, err := nc.Subscribe("order.created", func(msg *nats.Msg) {
		evt, err := events.UnmarshalOrderCreated(msg.Data)
		if err != nil {
			log.Println("error:", err)
			return
		}

		handler.Handle(evt)
	})
	if err != nil {
		log.Fatalf("Subscribe failed: %v", err)
	}

	log.Println("Inventory serivce running... listening to order.created")
	select {}
}
