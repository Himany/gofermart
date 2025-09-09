package handlers

import (
	"fmt"

	"github.com/Himany/gofermart/internal/models"
)

func validateRegisterJSON(r models.AuthDataRequest) error {
	if r.Login == "" {
		return fmt.Errorf("login is required")
	}
	if r.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func validateWithdrawJSON(r models.BonusWithdrawRequest) error {
	if r.Order == "" {
		return fmt.Errorf("order is required")
	}
	if r.Sum <= 0 {
		return fmt.Errorf("sum must be greater than zero")
	}
	return nil
}
