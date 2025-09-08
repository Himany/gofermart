package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	UserID int    `json:"uid"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}

func (h *Handler) signJWT(userID int, login string) (string, error) {
	ttl := h.TokenTTL
	if ttl == 0 {
		ttl = 12 * time.Hour
	}
	now := time.Now()
	claims := jwtClaims{
		UserID: userID,
		Login:  login,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.JWTSecret)
}

func (h *Handler) parseJWT(tokenStr string) (int, bool) {
	if len(h.JWTSecret) == 0 {
		return 0, false
	}
	tok, err := jwt.ParseWithClaims(
		tokenStr,
		&jwtClaims{},
		func(t *jwt.Token) (interface{}, error) { return h.JWTSecret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil || !tok.Valid {
		return 0, false
	}
	cl, ok := tok.Claims.(*jwtClaims)
	if !ok {
		return 0, false
	}
	return cl.UserID, true
}

func (h *Handler) setJWTCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(h.TokenTTL),
	})
}
func (h *Handler) setJWTHeader(w http.ResponseWriter, token string) {
	w.Header().Add("Authorization", "Bearer "+token)
}

func (h *Handler) authFromRequest(r *http.Request) (int, bool) {
	if c, err := r.Cookie("auth_token"); err == nil {
		if id, ok := h.parseJWT(c.Value); ok {
			return id, true
		}
	}
	auth := r.Header.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		if id, ok := h.parseJWT(auth[7:]); ok {
			return id, true
		}
	}
	return 0, false
}
