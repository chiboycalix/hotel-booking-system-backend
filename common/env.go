package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	prod := os.Getenv("PROD")
	if prod != "true" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}
	return nil
}
func EnvJWTSecret() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("JWT_SECRET")
}
func SenderEmail() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("SENDER_EMAIL")
}
func BrevoAPIKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("BREVO_API_KEY")
}
