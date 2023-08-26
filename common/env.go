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

func FrontendUrl() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("FRONTEND_URL")
}

func GoogleClientID() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("GOOGLE_CLIENT_ID")
}

func GoogleClientSecret() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("GOOGLE_CLIENT_SECRET")
}

func GoogleRedirectURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("GOOGLE_REDIRECT_URL")
}

func CloudinaryCloudName() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("CLOUDINARY_CLOUD_NAME")
}

func CloudinaryAPIKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("CLOUDINARY_API_KEY")
}

func CloudinaryAPISecret() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("CLOUDINARY_API_SECRET")
}

func CloudinaryUploadFolder() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
}
