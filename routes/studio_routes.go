package routes

import (
	"time"

	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func StudioRoutes(api fiber.Router) {
	studio := api.Group("/studio")

	//caching
	studioCache := cache.New(cache.Config{
		Expiration:   1 * time.Hour,
		CacheControl: true,
	})

	// Endpoint Publik (Bisa diakses user biasa tanpa login)
	studio.Get("/", studioCache, handlers.GetAllStudios)
	studio.Get("/:id", studioCache, handlers.GetStudioByID)

	// Endpoint Khusus Admin (Wajib login dan role = admin)
	studio.Post("/", middleware.JWTProtected, middleware.AdminOnly, handlers.CreateStudio)
	studio.Patch("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.UpdateStudio)
	studio.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteStudio)
}
