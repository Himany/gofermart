package models

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AddOrderRequest struct {
	OrderNumber string `json:"order_number"`
}

type BonusWithdrawRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}
