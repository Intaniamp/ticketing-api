package routes

import (
	"ticketing-api/internal/handlers"

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

	// --- CITY ---
	CityRoutes(api)

	// --- CINEMA ---
	CinemaRoutes(api)

	// --- STUDIO ---
	StudioRoutes(api)

	// --- SCHEDULE ---
	ScheduleRoutes(api)

	// --- SEAT ---
	SeatRoutes(api)
}
