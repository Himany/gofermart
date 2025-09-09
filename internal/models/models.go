package models

type AuthDataRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BonusWithdrawRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

type StatusOrder string

const (
	StatusRegistered StatusOrder = "REGISTERED"
	StatusInvalid    StatusOrder = "INVALID"
	StatusProcessing StatusOrder = "PROCESSING"
	StatusProcessed  StatusOrder = "PROCESSED"
)

type OrderInfo struct {
	Order   string      `json:"order"`
	Status  StatusOrder `json:"status"`
	Accrual *float64    `json:"accrual,omitempty"`
}
