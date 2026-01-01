package outbox

type Record struct {
	ID          int64
	AggregateID string
	EventType   string
	Payload     []byte
}
