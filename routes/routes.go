package routes

import (
	controllers "feedback-io.backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setups(app *fiber.App) {

	app.Get("/suggestions", controllers.GetSuggestions)
	app.Get("/suggestions/:id<int>", controllers.GetSuggestion)

	app.Put("/suggestions/:id<int>/vote", controllers.VoteSuggestion)

}
