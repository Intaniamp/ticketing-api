package routes

import (
	"ticketing-api/handlers"
	"ticketing-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func CinemaRoutes(api fiber.Router) {
	cinema := api.Group("/cinema")

	// Endpoint Publik (Bisa diakses user biasa tanpa login)
	cinema.Get("/", handlers.GetAllCinemas)
	cinema.Get("/:id", handlers.GetCinemaByID)

	// Endpoint Khusus Admin (Wajib login dan role = admin)
	cinema.Post("/", middleware.JWTProtected, middleware.AdminOnly, handlers.CreateCinema)
	cinema.Patch("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.UpdateCinema)
	cinema.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteCinema)
}
