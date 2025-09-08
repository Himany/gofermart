package server

import (
	"database/sql"
	"time"

	"github.com/Himany/gofermart/internal/config"
	"github.com/Himany/gofermart/internal/logger"
	"github.com/Himany/gofermart/internal/server/handlers"
	"github.com/Himany/gofermart/internal/storage"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Run(cfg *config.Config) error {
	var repo handlers.MarketRepo
	db, err := sql.Open("pgx", cfg.DataBaseURI)
	if err != nil {
		logger.Log.Error("failed to open database (postgres)", zap.Error(err))
	} else {
		repo, err = storage.NewPostgresStorage(db)
		if err != nil {
			logger.Log.Fatal("failed to init database storage", zap.Error(err))
		}
		defer db.Close()
	}

	handler := &handlers.Handler{
		Repo:      repo,
		JWTSecret: []byte(cfg.JWTSecret),
		TokenTTL:  time.Duration(cfg.JWTTokenTTL) * time.Second,
	}
	if err := Router(handler, cfg.RunAddress); err != nil {
		return err
	}

	return nil
}
