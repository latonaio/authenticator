package authorizer

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/form3tech-oss/jwt-go"
)

func generateJWT(userID int, exp int64) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = exp
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("PRIVATE_KEY")))
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}
	return token.SignedString(privateKey)
}

func TestValidateJWTToken(t *testing.T) {
	const privateKey = ""
	const publicKey = ""

	if err := os.Setenv("PRIVATE_KEY", privateKey); err != nil {
		t.Fatalf("failed to set env PRIVATE_KEY: %v", err)
	}
	if err := os.Setenv("AUTHENTICATOR_PUBLIC_KEY", publicKey); err != nil {
		t.Fatalf("failed to set env AUTHENTICATOR_PUBLIC_KEY: %v", err)
	}

	jwtStr, err := generateJWT(1, time.Now().Add(time.Minute).Unix())
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}

	expiredJwtStr, err := generateJWT(1, time.Now().Add(time.Minute*-1).Unix())
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}

	tests := []struct {
		tokenString string
		wantErr     bool
	}{
		{jwtStr, false},
		{jwtStr + "falsified suffix", true},
		{expiredJwtStr, true},
	}
	for i, test := range tests {
		_, err := VerifyJWTToken(test.tokenString)
		if (err != nil) != test.wantErr {
			t.Errorf("XXX ValidateJWTToken() error: %v, wantErr: %v", i, err, test.wantErr)
		}
	}
}
