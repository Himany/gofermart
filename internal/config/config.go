package config

import "go.uber.org/zap/zapcore"

type Config struct {
	RunAddress             string `env:"RUN_ADDRESS"`
	DataBaseURI            string `env:"DATABASE_URI"`
	ACCRUAL_SYSTEM_ADDRESS string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel               string `env:"LOGLEVEL"`
}

func (c *Config) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("runAddress", c.RunAddress)
	enc.AddString("dataBaseURI", c.DataBaseURI)
	enc.AddString("accrualSystemAddress", c.ACCRUAL_SYSTEM_ADDRESS)
	enc.AddString("logLevel", c.LogLevel)
	return nil
}
