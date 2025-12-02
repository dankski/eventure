package main

import (
	"log"
	"time"

	"eventure/libs/nats"
	"eventure/services/orders/internal/ports"
)

func main() {
	nc := nats.ConnectNATS()
	defer nc.Close()

	svc := ports.NewOrderService(nc)

	for {
		err := svc.CreateOrder("item-123", 1)
		if err != nil {
			log.Println("error:", err)
		}
		time.Sleep(5 * time.Second)
	}
}
