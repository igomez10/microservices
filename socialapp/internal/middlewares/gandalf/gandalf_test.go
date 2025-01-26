package gandalf

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	socialappjwt "github.com/igomez10/microservices/socialapp/internal/jwt"
)

func TestJWTParsing(t *testing.T) {
	hmacSampleSecret := []byte("secret")
	username := "ignacio"
	scopes := []string{"read", "write"}

	jwtToken := &socialappjwt.SocialAPPToken{
		Username:  username,
		Scopes:    scopes,
		Audience:  jwt.ClaimStrings{"aud1", "aud2"},
		Expires:   jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		NotBefore: jwt.NewNumericDate(time.Now().UTC().Add(-time.Minute)),
		Issuer:    "gandalf",
		Subject:   "social-app",
	}

	tokenString, err := socialappjwt.TokenToString(jwtToken, hmacSampleSecret)
	if err != nil {
		t.Fatal(err)
	}

	parsedToken, err := socialappjwt.ParseJWTToken(tokenString, hmacSampleSecret)
	if err != nil {
		t.Fatal(err)
	}

	if !parsedToken.Valid {
		t.Fatal("parsed token is invalid")
	}

	fromtoken, err := socialappjwt.FromToken(parsedToken)
	if err != nil {
		t.Fatal(err)
	}

	content1, err := json.Marshal(fromtoken)
	if err != nil {
		t.Fatal(err)
	}

	content2, err := json.Marshal(jwtToken)
	if err != nil {
		t.Fatal(err)
	}

	if string(content1) != string(content2) {
		fmt.Printf("%+v\n", content1)
		fmt.Printf("%+v\n", content2)
		t.Fatal("parsed token claims do not match")
	}

}
