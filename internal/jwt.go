package internal

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// ValidateJWT checks if the JWT is valid and if so returns the parsed token,
// otherwise it returns null end an error
func ValidateJWT(tokenString string) (*jwt.Token, error) {

	publicKeyString := GetString("JWT_PUBLIC_KEY", "")
	publicKeyPEM := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", publicKeyString)

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Check if the signing method is RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	return token, err
}
