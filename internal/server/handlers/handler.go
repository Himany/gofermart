package handlers

type MarketRepo interface {
	Ping() error
}

type Handler struct {
	Repo MarketRepo
	Key  string
}
