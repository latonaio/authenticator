package authorizer

import (
	"errors"
	"fmt"
	"os"

	"github.com/form3tech-oss/jwt-go"
)

// VerifyJWTToken verifies jwt token.
// Returns a token if the verified is successful, an error if it fails.
func VerifyJWTToken(tokenString string) (*jwt.Token, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("AUTHENTICATOR_PUBLIC_KEY")))
	if err != nil {
		return nil, fmt.Errorf("failed to parse authenticator public key: %v", err)
	}
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodRS256 {
			return nil, errors.New("invalid signing method")
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}
	if !parsedToken.Valid {
		return nil, errors.New("token is invalid")
	}
	return parsedToken, nil
}
