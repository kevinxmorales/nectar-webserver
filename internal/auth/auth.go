package auth

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/serialize"
)

type Cache interface {
	Get(key string) (value string, found bool, err error)
	Set(key, value string) error
}

type Store interface{}

type AuthenticationClient interface {
	VerifyIDToken(ctx context.Context, token string) (*AuthToken, error)
}

type Service struct {
	Store      Store
	AuthClient AuthenticationClient
	Cache      Cache
}

// NewService - returns a pointer to a new user service
func NewService(store Store, authClient AuthenticationClient, cache Cache) *Service {
	return &Service{
		Store:      store,
		AuthClient: authClient,
		Cache:      cache,
	}
}

func (s *Service) VerifyIDToken(ctx context.Context, sessionToken string) (*AuthToken, error) {
	serializedToken, found, err := s.Cache.Get(sessionToken)
	if err != nil {
		return nil, fmt.Errorf("failed to read from cache, this is not a key-not-found error: %v", err)
	}
	if !found {
		authToken, err := s.AuthClient.VerifyIDToken(ctx, sessionToken)
		if err != nil {
			return nil, fmt.Errorf("an error occurred verifying the auth token: %v", err)
		}
		serializedAuthToken, err := serializeToken(*authToken)
		if err != nil {
			return nil, fmt.Errorf("an error occurred serializing the auth token: %v", err)
		}
		if err := s.Cache.Set(sessionToken, serializedAuthToken); err != nil {
			return nil, fmt.Errorf("an error occurred setting the auth token in the cache: %v", err)
		}
	}
	authToken, err := getTokenFromSerializedForm(serializedToken)
	if err != nil {
		return nil, fmt.Errorf("invalid format of serialized token: %v", err)
	}
	return &authToken, nil
}

func serializeToken(authToken AuthToken) (string, error) {
	sx := serialize.SX{}
	sx["UID"] = authToken.UID
	sx["Claims"] = authToken.Claims
	return serialize.ToGOB64(sx)
}

func getTokenFromSerializedForm(serializedToken string) (AuthToken, error) {
	valMap, err := serialize.FromGOB64(serializedToken)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed to deserialize the authentication token: %v", err)
	}
	uid := valMap["UID"].(string)
	if uid == "" {
		return AuthToken{}, errors.New("no uid found in serialized token")
	}
	claims := valMap["Claims"].(map[string]interface{})
	if claims == nil {
		return AuthToken{}, errors.New("no claims found in serialized token")
	}
	authToken := AuthToken{
		UID:    uid,
		Claims: claims,
	}
	return authToken, nil
}
