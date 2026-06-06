package routes

import (
	"time"

	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func GenreRoutes(api fiber.Router) {
	genre := api.Group("/genre")

	//caching
	genreCache := cache.New(cache.Config{
		Expiration:   24 * time.Hour,
		CacheControl: true,
	})

	// Endpoint Publik (Bisa diakses user biasa tanpa login)
	genre.Get("/", genreCache, handlers.GetAllGenres)
	genre.Get("/:id", genreCache, handlers.GetGenreByID)

	// Endpoint Khusus Admin (Wajib login dan role = admin)
	genre.Post("/", middleware.JWTProtected, middleware.AdminOnly, handlers.CreateGenre)
	genre.Patch("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.UpdateGenre)
	genre.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteGenre)
}
