package events

type PaymentCharge struct {
	OrderID string `json:"order_id"`
}
