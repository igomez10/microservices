package jwt

import (
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type SocialAPPToken struct {
	Username  string           `json:"username"`
	Scopes    []string         `json:"scopes"`
	Audience  jwt.ClaimStrings `json:"aud"`
	Expires   *jwt.NumericDate `json:"exp"`
	IssuedAt  *jwt.NumericDate `json:"iat"`
	NotBefore *jwt.NumericDate `json:"nbf"`
	Issuer    string           `json:"iss"`
	Subject   string           `json:"sub"`
}

// Implement jwt.Claims interface methods
func (j *SocialAPPToken) GetAudience() (jwt.ClaimStrings, error) {
	return j.Audience, nil
}

func (j *SocialAPPToken) GetExpirationTime() (*jwt.NumericDate, error) {
	return j.Expires, nil
}

func (j *SocialAPPToken) GetIssuedAt() (*jwt.NumericDate, error) {
	return j.IssuedAt, nil
}

func (j *SocialAPPToken) GetIssuer() (string, error) {
	return j.Issuer, nil
}

func (j *SocialAPPToken) GetNotBefore() (*jwt.NumericDate, error) {
	return j.NotBefore, nil
}

func (j *SocialAPPToken) GetSubject() (string, error) {
	return j.Subject, nil
}

func ParseJWTToken(tokenString string, hmacSampleSecret []byte) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return parsedToken, nil
}

func TokenToString(token *SocialAPPToken, hmacSampleSecret []byte) (string, error) {
	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
	tokenString, err := authToken.SignedString(hmacSampleSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func FromToken(token *jwt.Token) (*SocialAPPToken, error) {
	newtoken := SocialAPPToken{}
	content, err := json.Marshal(token.Claims)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims: %w", err)
	}

	err = json.Unmarshal(content, &newtoken)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return &newtoken, nil
}

func NewTokenFromString(tokenString, secret string) (*SocialAPPToken, error) {
	jwtToken, err := ParseJWTToken(tokenString, []byte(secret))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return FromToken(jwtToken)
}
