package main

import (
	"log"
	"os"

	database "feedback-io.backend/config"
	"feedback-io.backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	})

	database.ConnectDatabase()

	routes.Setups(app)

	log.Fatal(app.Listen(":" + port))

}
