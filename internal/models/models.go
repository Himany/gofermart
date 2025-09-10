package models

type AuthDataRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BonusWithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
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

type OrderDTO struct {
	Number     string   `json:"number"`
	Status     string   `json:"status"`
	Accrual    *float64 `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"`
}

type WithdrawalDTO struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
