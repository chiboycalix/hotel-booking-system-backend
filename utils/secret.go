package utils

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/chiboycalix/hotel-booking-system-backend/common"
	"github.com/chiboycalix/hotel-booking-system-backend/responses"
)

type JWTClaim struct {
	ID      string `json:"_id"`
	IsAdmin bool   `json:"isAdmin"`
	jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func GenerateJWT(id string, isAdmin bool) (string, error) {
	claims := &JWTClaim{
		ID:      id,
		IsAdmin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	tokenSecret := common.EnvJWTSecret()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

func ValidateToken(tokenString string, secretKey string, c *fiber.Ctx) (JWTClaim, error) {
	var claims JWTClaim
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return JWTClaim{}, err
		}
		return JWTClaim{}, c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Unauthorized", Data: &fiber.Map{"error": err.Error()}})
	}
	if !token.Valid {
		return JWTClaim{}, c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Unauthorized", Data: &fiber.Map{"error": "Unauthorized"}})
	}
	return claims, nil
}

func IsAdmin(tokenString string, secretKey string, c *fiber.Ctx) (bool, error) {
	claims, err := ValidateToken(tokenString, secretKey, c)
	if err != nil {
		return false, err
	}
	return claims.IsAdmin, nil
}
