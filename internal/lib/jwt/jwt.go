package jwt

import (
	"fmt"
	"github.com/c1tad3l/wedo-auth-grpc-/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var Secret string = "eerstyrjndrfbsrvaegrthryj"

func NewTokenPair(user models.User, duration time.Duration) (at string, rt string, err error) {

	aToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": user.Uuid,
		"exp": duration,
	})

	accessToken, err := aToken.SignedString([]byte(Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %v", err)
	}
	rToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": user.Uuid,
	})

	refreshToken, err := rToken.SignedString([]byte(Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}
	return accessToken, refreshToken, nil
}
