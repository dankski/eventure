package events

import "encoding/json"

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func UnmarshalOrderCreated(data []byte) (*OrderCreated, error) {
	var evt OrderCreated
	err := json.Unmarshal(data, &evt)

	return &evt, err
}

func UnmarshalInventoryReserved(data []byte) (*InventoryReserved, error) {
	var evt InventoryReserved
	err := json.Unmarshal(data, &evt)

	return &evt, err
}

func UnmarshalInventoryFailed(data []byte) (*InventoryFailed, error) {
	var evt InventoryFailed
	err := json.Unmarshal(data, &evt)

	return &evt, err
}

func UnmarshalPaymentAuthorized(data []byte) (*PaymentAuthorized, error) {
	var evt PaymentAuthorized
	err := json.Unmarshal(data, &evt)

	return &evt, err
}

func UnmarshalPaymentCharge(data []byte) (*PaymentCharge, error) {
	var evt PaymentCharge
	err := json.Unmarshal(data, &evt)

	return &evt, err
}

func UnmarshalPaymentFailed(data []byte) (*PaymentFailed, error) {
	var evt PaymentFailed
	err := json.Unmarshal(data, &evt)

	return &evt, err
}
