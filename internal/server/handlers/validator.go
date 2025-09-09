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
