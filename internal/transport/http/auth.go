package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/auth"
	"net/http"
	"os"
	"strings"
)

type AuthService interface {
	Login(context.Context, string, string) (*auth.TokenDetails, error)
}

func JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
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
		if validateToken(authHeaderParts[1]) {
			original(w, r)
		} else {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}
	}
}

func validateToken(accessToken string) bool {
	var signingKey = []byte(os.Getenv("TOKEN_SECRET"))
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("could not validate auth token")
		}
		return signingKey, nil
	})
	if err != nil {
		log.Info("An error occurred while validating token")
		log.Error(err)
		return false
	}
	return token.Valid
}

type LoginRequest struct {
	Email    string `json:"email" required:"true"`
	Password string `json:"password" required:"true"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Error(err)
		http.Error(w, "unable to decode request", http.StatusInternalServerError)
		return
	}
	validate := validator.New()
	if err := validate.Struct(loginRequest); err != nil {
		log.Error(err)
		http.Error(w, "not a valid user object", http.StatusBadRequest)
		return
	}
	td, err := h.AuthService.Login(r.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		http.Error(w, "unable to authenticate", http.StatusUnauthorized)
		return
	}
	if err := json.NewEncoder(w).Encode(LoginResponse{Token: td.AccessToken, RefreshToken: td.RefreshToken}); err != nil {
		http.Error(w, "unable to encode JWT", http.StatusInternalServerError)
		return
	}
}
