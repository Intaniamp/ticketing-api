package routes

import (
	"ticketing-api/handlers"
	"ticketing-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func StudioRoutes(api fiber.Router) {
	studio := api.Group("/studio")

	// Endpoint Publik (Bisa diakses user biasa tanpa login)
	studio.Get("/", handlers.GetAllStudios)
	studio.Get("/:id", handlers.GetStudioByID)

	// Endpoint Khusus Admin (Wajib login dan role = admin)
	studio.Post("/", middleware.JWTProtected, middleware.AdminOnly, handlers.CreateStudio)
	studio.Patch("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.UpdateStudio)
	studio.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteStudio)
}
