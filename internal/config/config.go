package config

import "go.uber.org/zap/zapcore"

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DataBaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel             string `env:"LOGLEVEL"`
	JWTTokenTTL          int    `env:"JWT_TOKEN_TTL"`
	JWTSecret            string `env:"JWT_SECRET"`
}

func (c *Config) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("runAddress", c.RunAddress)
	enc.AddString("dataBaseURI", c.DataBaseURI)
	enc.AddString("accrualSystemAddress", c.AccrualSystemAddress)
	enc.AddString("logLevel", c.LogLevel)
	enc.AddInt("jwtTokenTTL", c.JWTTokenTTL)
	enc.AddString("jwtSecret", c.JWTSecret)
	return nil
}
