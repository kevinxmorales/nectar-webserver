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

type AuthClient struct {
	Client *auth.Client
}

type AuthToken struct {
	UID    string                 `json:"uid,omitempty"`
	Claims map[string]interface{} `json:"-"`
}

func NewAuthClient() (*AuthClient, error) {
	client, err := SetUpAuthClient()
	if err != nil {
		return nil, fmt.Errorf("an error occurred trying to set up auth client")
	}
	return &AuthClient{Client: client}, nil
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

func (ac *AuthClient) CreateUser(ctx context.Context, newUserId string, email string, password string) error {
	params := (&auth.UserToCreate{}).
		UID(newUserId).
		Email(email).
		EmailVerified(false).
		Password(password).
		Disabled(false)
	if _, err := ac.Client.CreateUser(ctx, params); err != nil {
		return err
	}
	return nil
}

func (ac *AuthClient) VerifyIDToken(ctx context.Context, token string) (*AuthToken, error) {
	t, err := ac.Client.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, err
	}
	authToken := &AuthToken{
		UID:    t.UID,
		Claims: t.Claims,
	}
	return authToken, nil
}
