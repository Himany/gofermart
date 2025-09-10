package utils

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SetStringIfUnset(envSet map[string]bool, envKey string, cfgValue *string, flagValue string) {
	if envSet[envKey] {
		return
	}
	*cfgValue = flagValue
}

func SetIntIfUnset(envSet map[string]bool, envKey string, cfgValue *int, flagValue int) {
	if envSet[envKey] {
		return
	}
	*cfgValue = flagValue
}

func ParseRetryAfter(v string) time.Duration {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	if sec, err := strconv.Atoi(v); err == nil && sec >= 0 {
		return time.Duration(sec) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		return time.Until(t)
	}
	return 0
}
