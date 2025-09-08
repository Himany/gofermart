package handlers

import (
	"time"
)

type MarketRepo interface {
	Ping() error
	CreateUser(login, passwordHash string) (int, error)
	GetUserByLogin(login string) (int, string, error)
}

type Handler struct {
	Repo      MarketRepo
	JWTSecret []byte
	TokenTTL  time.Duration
}
