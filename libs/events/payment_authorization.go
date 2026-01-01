package events

type PaymentAuthorized struct {
	OrderID string `json:"order_id"`
}
