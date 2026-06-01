package routes

import (
	"ticketing-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// --- AUTH ---
	api.Post("/login", handlers.Login)
	api.Post("/register", handlers.Register)

	// --- USER ---
	UserRoutes(api)

	// --- FILM ---
	FilmRoutes(api)

	// --- GENRE ---
	GenreRoutes(api)
}
