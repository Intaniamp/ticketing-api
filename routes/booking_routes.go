package routes

import (
	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func BookingRoutes(api fiber.Router) {
	booking := api.Group("/booking")

	// Semua fitur booking wajib login (JWT Protected)
	booking.Post("/", middleware.JWTProtected, handlers.CreateBooking)
	booking.Get("/my-history", middleware.JWTProtected, handlers.GetMyBookings)
	booking.Get("/:id", middleware.JWTProtected, handlers.GetBookingByID)
	booking.Get("/", middleware.JWTProtected, middleware.AdminOnly, handlers.GetAllBookings) // Admin only
}
