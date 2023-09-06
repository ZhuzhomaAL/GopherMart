package clients

type LoyalClient interface {
	GetOrderProcessingInfo(order string) (OrderLoyaltyInfo, error)
}

const (
	StatusRegistered string = "REGISTERED"
	StatusProcessing string = "PROCESSING"
	StatusInvalid    string = "INVALID"
	StatusProcessed  string = "PROCESSED"
)

type OrderLoyaltyInfo struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
