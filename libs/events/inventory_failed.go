package events

type InventoryFailed struct {
	OrderID string `json:"order_id"`
	Reason  string `json:"reason"`
}
