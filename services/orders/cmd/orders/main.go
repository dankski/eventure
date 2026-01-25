package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	// Start saga when order is created
	startOrder := func(orderID []byte) {
		log.Println("OrderSaga: start")
		nc.Publish("orders.reserve_inventory", orderID)
	}

	// Inventory success
	nc.Subscribe("inventory.reserved", func(msg *nats.Msg) {
		log.Println("OrderSaga: inventory reserved")
		nc.Publish("orders.process_payment", msg.Data)
	})

	// Inventory failure
	nc.Subscribe("inventory.failed", func(msg *nats.Msg) {
		log.Println("OrderSaga: inventory failed → cancel order")
	})

	// Payment success
	nc.Subscribe("payment.succeeded", func(msg *nats.Msg) {
		log.Println("OrderSaga: payment succeeded → order completed")
	})

	// Payment failure
	nc.Subscribe("payment.failed", func(msg *nats.Msg) {
		log.Println("OrderSaga: payment failed → compensate inventory")
		nc.Publish("orders.compensate_inventory", msg.Data)
	})

	// Kick off demo order
	startOrder([]byte("order-123"))

	select {}
}
