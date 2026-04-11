package routes

import (
	"ticketing-api/handlers"
	"ticketing-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(api fiber.Router) {
	user := api.Group("/user")
	user.Get("/", middleware.JWTProtected, middleware.AdminOnly, handlers.GetAllUsers)
	user.Get("/:id", middleware.JWTProtected, handlers.GetUserByID)
	user.Patch("/:id", middleware.JWTProtected, handlers.UpdateUser)
	user.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteUser)
}
