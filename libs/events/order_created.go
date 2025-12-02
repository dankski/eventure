package events

type OrderCreated struct {
	OrderID string `json:"order_id"`
	ItemID  string `json:"item_id"`
	Qty     int    `json:"qty"`
}
