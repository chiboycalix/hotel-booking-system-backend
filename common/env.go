package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	// check if prod
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
