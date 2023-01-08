package http

import (
	"context"
	firebaseAuth "firebase.google.com/go/v4/auth"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type AuthService interface {
	VerifyIDToken(ctx context.Context, token string) (*firebaseAuth.Token, error)
}

func (h *Handler) JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if true {
			original(w, r)
			return
		}
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
		log.Info(token.UID)
		original(w, r)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	return
}
