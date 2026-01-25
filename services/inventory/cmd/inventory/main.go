package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	fmt.Println("Connected to NATS server at", nats.DefaultURL)

	nc.Subscribe("orders.reserve_inventory", func(msg *nats.Msg) {
		log.Println("Invetory: reserved stock")

		nc.Publish("inventory.reserved", msg.Data)
	})

	nc.Subscribe("orders.compensate_inventory", func(msg *nats.Msg) {
		log.Println("Inventory: compensation executed")
	})

	select {}
}
