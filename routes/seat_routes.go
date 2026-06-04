package routes

import (
	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SeatRoutes(api fiber.Router) {
	seat := api.Group("/seat")

	// Endpoint Publik (Bisa diakses user biasa untuk memilih kursi saat booking)
	seat.Get("/studio/:studio_id", handlers.GetSeatsByStudio)

	// Endpoint Khusus Admin (Wajib login dan role = admin)
	seat.Post("/bulk", middleware.JWTProtected, middleware.AdminOnly, handlers.BulkCreateSeats)
	seat.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteSeat)
}
