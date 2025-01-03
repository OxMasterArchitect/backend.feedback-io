package routes

import (
	controllers "feedback-io.backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setups(app *fiber.App) {

	app.Get("/suggestions", controllers.GetSuggestions)
	app.Get("/suggestions/:id", controllers.GetSuggestion)

	app.Put("/suggestions/:id/vote", controllers.VoteSuggestion)

}
