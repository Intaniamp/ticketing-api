package routes

import (
	"time"

	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func CinemaRoutes(api fiber.Router) {
	cinema := api.Group("/cinema")

	//caching
	cinemaCache := cache.New(cache.Config{
		Expiration:   24 * time.Hour,
		CacheControl: true,
	})

	// Endpoint Publik (Bisa diakses user biasa tanpa login)
	cinema.Get("/", cinemaCache, handlers.GetAllCinemas)
	cinema.Get("/:id", cinemaCache, handlers.GetCinemaByID)

	// Endpoint Khusus Admin (Wajib login dan role = admin)
	cinema.Post("/", middleware.JWTProtected, middleware.AdminOnly, handlers.CreateCinema)
	cinema.Patch("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.UpdateCinema)
	cinema.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteCinema)
}
