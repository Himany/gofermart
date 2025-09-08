package main

import (
	"flag"

	"github.com/Himany/gofermart/internal/config"
	"github.com/Himany/gofermart/internal/utils"
	"github.com/caarlos0/env/v11"
)

var envSet = map[string]bool{}

// Стандратные значения
const defaultRunAddress = "localhost:8080"
const defaultDataBaseURI = "host=localhost user=postgres password=123321 dbname=gofermart sslmode=disable"
const defaultAccrualSystemAddress = "localhost:8081"
const defaultLogLevel = "info"
const defaultJWTTokenTTL = 100
const defaultJWTSecret = "mysecret"

func parseFlags() (*config.Config, error) {
	var flagRunAddress = flag.String("a", defaultRunAddress, "address and port to run server")
	var flagDataBaseURI = flag.String("d", defaultDataBaseURI, "a string with settings for connecting the postgresql database")
	var flagAccrualSystemAddress = flag.String("f", defaultAccrualSystemAddress, "address of the accrual calculation system")
	var flagLogLevel = flag.String("l", defaultLogLevel, "log level")
	var flagJWTTokenTTL = flag.Int("t", defaultJWTTokenTTL, "JWT token TTL in seconds")
	var flagJWTSecret = flag.String("s", defaultJWTSecret, "JWT secret")

	flag.Parse()

	cfg := config.Config{
		RunAddress:           *flagRunAddress,
		DataBaseURI:          *flagDataBaseURI,
		AccrualSystemAddress: *flagAccrualSystemAddress,
		LogLevel:             *flagLogLevel,
		JWTTokenTTL:          *flagJWTTokenTTL,
		JWTSecret:            *flagJWTSecret,
	}

	err := env.ParseWithOptions(&cfg, env.Options{
		OnSet: func(tag string, value any, isDefault bool) {
			if !isDefault {
				envSet[tag] = true
			}
		},
	})
	if err != nil {
		return nil, err
	}

	utils.SetStringIfUnset(envSet, "RUN_ADDRESS", &cfg.RunAddress, *flagRunAddress)
	utils.SetStringIfUnset(envSet, "DATABASE_URI", &cfg.DataBaseURI, *flagDataBaseURI)
	utils.SetStringIfUnset(envSet, "ACCRUAL_SYSTEM_ADDRESS", &cfg.AccrualSystemAddress, *flagAccrualSystemAddress)
	utils.SetStringIfUnset(envSet, "LOGLEVEL", &cfg.LogLevel, *flagLogLevel)
	utils.SetIntIfUnset(envSet, "JWT_TOKEN_TTL", &cfg.JWTTokenTTL, *flagJWTTokenTTL)
	utils.SetStringIfUnset(envSet, "JWT_SECRET", &cfg.JWTSecret, *flagJWTSecret)

	return &cfg, nil
}
