package authorizer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/form3tech-oss/jwt-go"
)

// VerifyJWTToken verifies jwt token.
// Returns a token if the verified is successful, an error if it fails.
func VerifyJWTToken(tokenString string) (*jwt.Token, error) {
	b,err := ioutil.ReadFile(os.Getenv("CREDENTIAL_FILE_PATH"))
	if err != nil {
		return nil,fmt.Errorf("failed to read authenticator public key file: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(b)
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
