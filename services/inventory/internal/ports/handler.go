package ports

import (
	"log"

	"eventure/libs/events"
)

type InventoryHandler struct{}

func NewInventoryHandler() *InventoryHandler {
	return &InventoryHandler{}
}

func (h *InventoryHandler) Handle(evt *events.OrderCreated) {
	log.Printf("InventoryService: received order.created for order %s, item=%s qty=%d\n", evt.OrderID, evt.ItemID, evt.Qty)
}
