package main

import (
	"log"
	"os"

	"github.com/chiboycalix/hotel-booking-system-backend/common"
	"github.com/chiboycalix/hotel-booking-system-backend/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	// init env
	err := common.LoadEnv()
	if err != nil {
		return err
	}

	// init db
	err = common.InitDB()
	if err != nil {
		return err
	}

	// defer closing db
	defer common.CloseDB()

	// create app
	app := fiber.New()
	allowOriginsFunc := func(origin string) bool {
		// Replace this logic with your own rules
		allowedDomains := []string{
			"https://localhost:3000",
			"http://localhost:3000",
			"https://localhost:3001",
			"http://localhost:3001",
		}

		for _, domain := range allowedDomains {
			if origin == domain {
				return true
			}
		}
		return false
	}

	// add basic middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Access-Control-Allow-Credentials",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowOriginsFunc: allowOriginsFunc,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	app.Use(recover.New())
	// routes
	router.UserRoute(app)
	router.AuthRoutes(app)
	router.ListingRoutes(app)
	// start server
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "3001"
	}
	log.Fatal(app.Listen(":" + port))
	return nil
}
