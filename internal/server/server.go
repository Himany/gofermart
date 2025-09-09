package server

import (
	"database/sql"
	"time"

	"github.com/Himany/gofermart/internal/accrual"
	"github.com/Himany/gofermart/internal/config"
	"github.com/Himany/gofermart/internal/logger"
	"github.com/Himany/gofermart/internal/server/handlers"
	"github.com/Himany/gofermart/internal/storage"
	"github.com/Himany/gofermart/internal/worker"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Run(cfg *config.Config) error {
	db, err := sql.Open("pgx", cfg.DataBaseURI)
	if err != nil {
		logger.Log.Fatal("failed to open database (postgres)", zap.Error(err))
	}

	if err := db.Ping(); err != nil {
		logger.Log.Fatal("failed to ping database", zap.Error(err))
	}

	repo, err := storage.NewPostgresStorage(db)
	if err != nil {
		logger.Log.Fatal("failed to init database storage", zap.Error(err))
	}
	defer db.Close()

	accrual := accrual.New(cfg.AccrualSystemAddress, 5*time.Second)
	pool := worker.NewPool(repo, accrual, 4, 1024)
	pool.Start()
	defer pool.Stop()

	handler := &handlers.Handler{
		Repo:        repo,
		JWTSecret:   []byte(cfg.JWTSecret),
		TokenTTL:    time.Duration(cfg.JWTTokenTTL) * time.Second,
		AccrualPool: pool,
	}
	if err := Router(handler, cfg.RunAddress); err != nil {
		return err
	}

	return nil
}
