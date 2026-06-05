package routes

import (
	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func PaymentRoutes(api fiber.Router) {
	payment := api.Group("/payment")

	// Endpoint Khusus Admin (Bisa lihat semua pembayaran masuk)
	payment.Get("/", middleware.JWTProtected, middleware.AdminOnly, handlers.GetAllPayments)

	// Endpoint User Biasa
	payment.Post("/", middleware.JWTProtected, handlers.ProcessPayment)
	payment.Get("/my-history", middleware.JWTProtected, handlers.GetMyPayments)
}
