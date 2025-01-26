package authentication

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWT(t *testing.T) {
	hmacSampleSecret := []byte("secret")
	username := "ignacio"
	scopes := []string{"read", "write"}
	validUntil := time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"nbf":      time.Now().Unix(),
		"exp":      validUntil,
		"scopes":   scopes,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(tokenString)
}
