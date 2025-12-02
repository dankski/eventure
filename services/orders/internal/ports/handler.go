package ports

import (
	"log"

	"eventure/libs/events"

	. "eventure/services/orders/internal/domain/order"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type OrderService struct {
	nc    *nats.Conn
	store *Store
}

func NewOrderService(nc *nats.Conn, store *Store) *OrderService {
	return &OrderService{nc: nc, store: store}
}

func (o *OrderService) CreateOrder(itemID string, qty int) (string, error) {
	orderID := uuid.NewString()
	o.store.SetStatus(orderID, StatusNew)

	evt := events.OrderCreated{
		OrderID: orderID,
		ItemID:  itemID,
		Qty:     qty,
	}

	data, _ := events.Marshal(evt)
	log.Printf("OrderService: created order %s", orderID)
	o.nc.Publish("order.created", data)

	return orderID, nil
}

func (o *OrderService) StartSagaListeners() {

	o.nc.Subscribe("inventory.reserved", func(msg *nats.Msg) {
		evt, _ := events.UnmarshalInventoryReserved(msg.Data)
		log.Printf("OrderService: invetory for order %s", evt.OrderID)
		o.store.SetStatus(evt.OrderID, StatusReserved)

		evtPaymentCharge := events.PaymentCharge{OrderID: evt.OrderID}
		data, _ := events.Marshal(evtPaymentCharge)
		o.nc.Publish("payment.charge", data) // trigger payment service
	})

	// inventory failed -> cancel
	o.nc.Subscribe("inventory.failed", func(msg *nats.Msg) {
		evt, _ := events.UnmarshalInventoryFailed(msg.Data)
		log.Printf("OrderService: invetory faield for order: %s: %s", evt.OrderID, evt.Reason)
		o.store.SetStatus(evt.OrderID, StatusCancelled)
	})

	// payment authorized -> complete order
	o.nc.Subscribe("payment.authorize", func(msg *nats.Msg) {
		evt, _ := events.UnmarshalPaymentAuthorized(msg.Data)
		log.Printf("OrderService: payment authorized for order %s", evt.OrderID)
		o.store.SetStatus(evt.OrderID, StatusCompleted)
	})

	// payment failed -> cancel
	o.nc.Subscribe("payment.failed", func(msg *nats.Msg) {
		evt, _ := events.UnmarshalPaymentFailed(msg.Data)
		log.Printf("OrderService: payment failed for order %s: %s", evt.OrderID, evt.Reason)
		o.store.SetStatus(evt.OrderID, StatusCancelled)
	})
}
