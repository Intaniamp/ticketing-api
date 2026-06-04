package routes

import (
	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func FilmRoutes(api fiber.Router) {
	film := api.Group("/film")
	film.Get("/", handlers.GetAllFilm)
	film.Get("/:id", handlers.GetFilmByID)

	// Only admin
	film.Post("/", middleware.JWTProtected, middleware.AdminOnly, handlers.CreateFilm)
	film.Patch("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.UpdateFilm)
	film.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteFilm)
	film.Post("/:id/poster", middleware.JWTProtected, middleware.AdminOnly, handlers.UploadPoster)
}
