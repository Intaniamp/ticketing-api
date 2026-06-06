package routes

import (
	"time"

	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func FilmRoutes(api fiber.Router) {
	film := api.Group("/film")

	//caching
	filmCache := cache.New(cache.Config{
		Expiration:   30 * time.Minute,
		CacheControl: true,
	})

	// Endpoint Publik (Bisa diakses user biasa tanpa login)
	film.Get("/", filmCache, handlers.GetAllFilm)
	film.Get("/:id", filmCache, handlers.GetFilmByID)

	// Only admin
	film.Post("/", middleware.JWTProtected, middleware.AdminOnly, handlers.CreateFilm)
	film.Patch("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.UpdateFilm)
	film.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteFilm)
	film.Post("/:id/poster", middleware.JWTProtected, middleware.AdminOnly, handlers.UploadPoster)
}
