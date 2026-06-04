package routes

import (
	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func CityRoutes(api fiber.Router) {
	city := api.Group("/city")

	// Endpoint Publik (Bisa diakses user biasa tanpa login)
	city.Get("/", handlers.GetAllCities)
	city.Get("/:id", handlers.GetCityByID)

	// Endpoint Khusus Admin (Wajib login dan role = admin)
	city.Post("/", middleware.JWTProtected, middleware.AdminOnly, handlers.CreateCity)
	city.Patch("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.UpdateCity)
	city.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteCity)
}
