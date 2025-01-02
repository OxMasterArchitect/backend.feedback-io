package main

import (
	"log"

	database "feedback-io.backend/config"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	})

	database.ConnectDatabase()

	log.Fatal(app.Listen(":8000"))

}
