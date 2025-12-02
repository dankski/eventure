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
