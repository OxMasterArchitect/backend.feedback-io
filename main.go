package main

import (
	"log"

	database "feedback-io.backend/config"
	"feedback-io.backend/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	})

	database.ConnectDatabase()

	routes.Setups(app)

	log.Fatal(app.Listen(":80"))

}
