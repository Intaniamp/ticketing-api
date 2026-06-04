package routes

import (
	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func GenreRoutes(api fiber.Router) {
	genre := api.Group("/genre")

	// Endpoint Publik (Bisa diakses user biasa tanpa login)
	genre.Get("/", handlers.GetAllGenres)
	genre.Get("/:id", handlers.GetGenreByID)

	// Endpoint Khusus Admin (Wajib login dan role = admin)
	genre.Post("/", middleware.JWTProtected, middleware.AdminOnly, handlers.CreateGenre)
	genre.Patch("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.UpdateGenre)
	genre.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteGenre)
}
