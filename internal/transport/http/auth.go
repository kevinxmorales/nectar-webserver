package http

import (
	"context"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/auth"
	"net/http"
	"strings"
)

type AuthService interface {
	VerifyIDToken(ctx context.Context, token string) (*auth.AuthToken, error)
}

func (h *Handler) JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}
		// Bearer token-string
		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}
		token, err := h.AuthService.VerifyIDToken(r.Context(), authHeaderParts[1])
		if err != nil {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}
		userId := token.Claims["user_id"]
		newCtx := context.WithValue(r.Context(), "userId", userId)
		reqWithContext := r.WithContext(newCtx)
		original(w, reqWithContext)
	}
}
