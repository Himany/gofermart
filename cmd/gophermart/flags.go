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
const defaultAccrualSystemAddress = "localhost:8080"
const defaultLogLevel = "info"

func parseFlags() (*config.Config, error) {
	var flagRunAddress = flag.String("a", defaultRunAddress, "address and port to run server")
	var flagDataBaseURI = flag.String("d", defaultDataBaseURI, "a string with settings for connecting the postgresql database")
	var flagAccrualSystemAddress = flag.String("f", defaultAccrualSystemAddress, "address of the accrual calculation system")
	var flagLogLevel = flag.String("l", defaultLogLevel, "log level")

	flag.Parse()

	var cfg config.Config
	err := env.ParseWithOptions(&cfg, env.Options{
		OnSet: func(tag string, value any, isDefault bool) {
			envSet[tag] = true
		},
	})
	if err != nil {
		return nil, err
	}

	utils.SetStringIfUnset(envSet, "RUN_ADDRESS", &cfg.RunAddress, *flagRunAddress)
	utils.SetStringIfUnset(envSet, "DATABASE_URI", &cfg.DataBaseURI, *flagDataBaseURI)
	utils.SetStringIfUnset(envSet, "ACCRUAL_SYSTEM_ADDRESS", &cfg.ACCRUAL_SYSTEM_ADDRESS, *flagAccrualSystemAddress)
	utils.SetStringIfUnset(envSet, "LOGLEVEL", &cfg.LogLevel, *flagLogLevel)

	return &cfg, nil
}
