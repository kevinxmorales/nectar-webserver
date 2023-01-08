package auth

import (
	"context"
	firebase "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
	"fmt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"google.golang.org/api/option"
)

func SetUpAuthClient() (*firebaseAuth.Client, error) {
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase auth: %v", err)
	}
	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase auth: %v", err)
	}
	return authClient, nil
}

type Store interface {
	GetCredentialsByEmail(context.Context, string) (user.User, error)
}

type Service struct {
	Store      Store
	AuthClient *firebaseAuth.Client
}

// NewService - returns a pointer to a new user service
func NewService(store Store, authClient *firebaseAuth.Client) *Service {
	return &Service{
		Store:      store,
		AuthClient: authClient,
	}
}

func (s *Service) VerifyIDToken(ctx context.Context, token string) (*firebaseAuth.Token, error) {
	return s.AuthClient.VerifyIDToken(ctx, token)
}
