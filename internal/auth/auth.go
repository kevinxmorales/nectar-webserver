package auth

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"os"
)

// decrypt from base64 to decrypted string
func decrypt(keyString string, stringToDecrypt string) error {
	ciphertext, _ := base64.URLEncoding.DecodeString(stringToDecrypt)

	block, err := aes.NewCipher([]byte(keyString))
	if err != nil {
		return err
	}

	if len(ciphertext) < aes.BlockSize {
		return fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	decryptedText := fmt.Sprintf("%s", ciphertext)

	// create a new file for saving the encrypted data.
	f, err := os.Create("decryptedServiceAccountKey.json")
	if err != nil {
		return err
	}
	if _, err = io.Copy(f, bytes.NewReader([]byte(decryptedText))); err != nil {
		return err
	}
	return nil
}

func SetUpAuthClient() (*auth.Client, error) {
	secret := os.Getenv("ENCRYPT_SECRET")
	plaintext, err := ioutil.ReadFile("encryptedServiceAccountKey.json")
	if err != nil {
		return nil, fmt.Errorf("error initializing auth: %s", err)
	}
	if err := decrypt(secret, string(plaintext)); err != nil {
		return nil, err
	}
	opt := option.WithCredentialsFile("decryptedServiceAccountKey.json")
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
	AuthClient *auth.Client
}

// NewService - returns a pointer to a new user service
func NewService(store Store, authClient *auth.Client) *Service {
	return &Service{
		Store:      store,
		AuthClient: authClient,
	}
}

func (s *Service) VerifyIDToken(ctx context.Context, token string) (*auth.Token, error) {
	return s.AuthClient.VerifyIDToken(ctx, token)
}
