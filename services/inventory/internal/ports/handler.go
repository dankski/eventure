package ports

import (
	"eventure/libs/events"
	"log"
	"math/rand"

	"github.com/nats-io/nats.go"
)

type InventoryHandler struct {
	nc *nats.Conn
}

func NewInventoryHandler(nc *nats.Conn) *InventoryHandler {
	return &InventoryHandler{nc: nc}
}

func (h *InventoryHandler) Handle(evt *events.OrderCreated) {
	// simulate random inventory outcome
	if rand.Intn(100) < 80 { // 80% chance of success
		ack := &events.InventoryReserved{
			OrderID: evt.OrderID,
		}
		data, _ := events.Marshal(ack)
		log.Printf("InventoryService: reserved inventory for order %s", evt.OrderID)
		h.nc.Publish("inventory.reserved", data)
	} else {
		fail := events.InventoryFailed{OrderID: evt.OrderID, Reason: "out of stock"}
		data, _ := events.Marshal(fail)
		log.Printf("InventoryService: failed to reserve inventory for order %s", evt.OrderID)
		h.nc.Publish("inventory.failed", data)
	}
}
