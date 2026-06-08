package routes

import (
	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func SeatRoutes(api fiber.Router) {
	seat := api.Group("/seat")
	//caching
	seatCache := cache.New(cache.Config{
		Expiration:   1 * time.Hour,
		CacheControl: true,
	})

	// 1. Ambil data master kursi fisik di dalam studio (biasanya untuk admin atau cek denah dasar), jarang berubah.
	seat.Get("/studio/:studio_id", seatCache, handlers.GetSeatsByStudio)

	// 2. Ambil data status kursi LIVE (available/booked) berdasarkan jadwal tayang film! (Ini yang dipakai FE), 
	// ini sering berubah karena berdasarkan pilihan jadwal tayang, apakah seat booked/avail.
	seat.Get("/schedule/:schedule_id", handlers.GetSeatsBySchedule)

	// Endpoint Khusus Admin
	seat.Post("/bulk", middleware.JWTProtected, middleware.AdminOnly, handlers.BulkCreateSeats)
	seat.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteSeat)
	seat.Delete("/studio/:studio_id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteSeatsByStudio)
}
