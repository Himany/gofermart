package handlers

import (
	"time"

	"github.com/Himany/gofermart/internal/worker"
)

type MarketRepo interface {
	Ping() error

	CreateUser(login, passwordHash string) (int, error)
	GetUserByLogin(login string) (int, string, error)

	GetOrderOwner(number string) (userID int, found bool, err error)
	AddOrder(userID int, number string) error
	UpdateOrderStatus(number string, status string, accrual *float64) error
	AddAccrual(userID int, orderNumber string, amount float64) error

	GetBalance(userID int) (float64, error)
	AddWithdraw(userID int, orderNumber string, amount float64) error
}

type Handler struct {
	Repo        MarketRepo
	JWTSecret   []byte
	TokenTTL    time.Duration
	AccrualPool *worker.Pool
}
