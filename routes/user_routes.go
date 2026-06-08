package routes

import (
	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(api fiber.Router) {
	user := api.Group("/user")
	user.Get("/:id", middleware.JWTProtected, handlers.GetUserByID)
	user.Patch("/:id", middleware.JWTProtected, handlers.UpdateUser)

	// Only admin, divalidate di middleware, bukan di handler, karena handler bisa dipakai untuk endpoint lain yang tidak harus admin.
	user.Get("/", middleware.JWTProtected, middleware.AdminOnly, handlers.GetAllUsers)
	user.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteUser) 
}
