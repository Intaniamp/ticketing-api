package routes

import (
	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RefundRoutes(api fiber.Router) {
	refund := api.Group("/refund")

	// Wajib login untuk melakukan refund
	refund.Post("/", middleware.JWTProtected, handlers.ProcessRefund)
}
