package utils

import (
	"time"

	"github.com/chiboycalix/hotel-booking-system-backend/common"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaim struct {
	ID string `json:"_id"`
	jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func GenerateJWT(id string) (string, error) {
	claims := &JWTClaim{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	tokenSecret := common.EnvJWTSecret()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}
