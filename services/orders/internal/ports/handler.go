package ports

import (
	"log"

	"eventure/libs/events"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type OrderService struct {
	nc *nats.Conn
}

func NewOrderService(nc *nats.Conn) *OrderService {
	return &OrderService{nc: nc}
}

func (o *OrderService) CreateOrder(itemID string, qty int) error {
	evt := events.OrderCreated{
		OrderID: uuid.NewString(),
		ItemID:  itemID,
		Qty:     qty,
	}

	data, err := events.Marshal(evt)
	if err != nil {
		return err
	}

	log.Printf("OrderService: publishing event order.created: %+v\n", evt)

	return o.nc.Publish("order.created", data)
}
