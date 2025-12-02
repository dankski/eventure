package ports

import (
	"eventure/libs/events"
	"log"
	"math/rand"

	"github.com/nats-io/nats.go"
)

type PaymentHandler struct {
	nc *nats.Conn
}

func NewPaymentHandler(nc *nats.Conn) *PaymentHandler {
	return &PaymentHandler{nc: nc}
}

func (h *PaymentHandler) Charge(orderID string) {
	if rand.Intn(100) < 90 { // 90% chance of success
		ack := &events.PaymentAuthorized{
			OrderID: orderID,
		}
		data, _ := events.Marshal(ack)
		log.Printf("PaymentService: authorized payment for order %s", orderID)
		h.nc.Publish("payment.authorize", data)
	} else {
		fail := &events.PaymentFailed{
			OrderID: orderID,
			Reason:  "card declined",
		}
		data, _ := events.Marshal(fail)
		log.Printf("PaymentService: failed to authorize payment for order %s", orderID)
		h.nc.Publish("payment.failed", data)
	}
}
