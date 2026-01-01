package main

import (
	"fmt"
	"log"
	"time"

	"eventure/libs/nats"
	"eventure/services/orders/internal/domain/order"
	"eventure/services/orders/internal/ports"
)

func main() {
	nc := nats.ConnectNATS()
	defer nc.Close()

	svc := ports.NewOrderService(nc, order.NewStore())

	svc.StartSagaListeners()

	for {
		val, err := svc.CreateOrder("item-123", 1)
		if err != nil {
			log.Println("error:", err)
		}

		fmt.Println("Created order:", val)
		time.Sleep(5 * time.Second)
	}
}
