package jwt

import (
	"github.com/c1tad3l/wedo-auth-grpc-/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var Secret string = "eerstyrjndrfbsrvaegrthryj"

func NewToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_uuid"] = user.Uuid
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
