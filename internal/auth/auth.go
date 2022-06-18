package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

// TokenDetails For storing in Redis
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
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

func CreateToken(usr user.User) (*TokenDetails, error) {
	td := TokenDetails{}
	expirationTime := time.Now().Add((((1 * time.Hour) * 24) * 7) * 52)
	claims := &Claims{
		ID:        usr.ID,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Email:     usr.Email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	accessTokenSigned, err := accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}
	expirationTime = time.Now().Add(time.Hour * 24)
	refreshClaims := RefreshClaims{
		ID: usr.ID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenSigned, err := refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	td.AccessToken = accessTokenSigned
	td.RefreshToken = refreshTokenSigned
	return &td, nil
}

// Check if two passwords do not match using Bcrypt's CompareHashAndPassword
// return nil on success and an error on failure
func passwordsDoNotMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currPassword))
	return err != nil
}

func (s *Service) Login(ctx context.Context, email string, givenPassword string) (*TokenDetails, error) {
	// Get the expected password from our in memory map
	usr, err := s.Store.GetCredentialsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if passwordsDoNotMatch(usr.Password, givenPassword) {
		return nil, fmt.Errorf("unauthorized")
	}
	return CreateToken(usr)
}
