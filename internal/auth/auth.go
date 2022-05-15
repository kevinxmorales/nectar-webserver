package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"time"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

type Store interface {
	GetCredentialsByEmail(context.Context, string) (user.User, error)
}

type Service struct {
	Store Store
}

// NewService - returns a pointer to a new user service
func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) Login(ctx context.Context, email string, givenPassword string) (string, error) {
	// Get the expected password from our in memory map
	usr, err := s.Store.GetCredentialsByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	// If a password exists for the given user
	// AND, if it is the same as the password we received, then we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if usr.Password != givenPassword {
		return "", fmt.Errorf("unauthorized")
	}
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		ID:    usr.ID,
		Name:  usr.Name,
		Email: usr.Email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte("nectar"))
	if err != nil {
		log.Error(err)
		return "", err
	}
	return tokenString, nil
}
