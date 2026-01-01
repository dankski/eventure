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

	handler := ports.NewInventoryHandler(nc)

	_, _ = nc.Subscribe("order.created", func(msg *nats.Msg) {
		evt, _ := events.UnmarshalOrderCreated(msg.Data)
		handler.Handle(evt)
	})

	log.Println("Inventory serivce running...")
	select {}
}
