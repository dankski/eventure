package order

type Status string

const (
	StatusNew       Status = "NEW"
	StatusReserved  Status = "RESERVED"
	StatusCharged   Status = "CHARGED"
	StatusCompleted Status = "COMPLETED"
	StatusCancelled Status = "CANCELLED"
)
