package routes

import (
	"ticketing-api/internal/handlers"
	"ticketing-api/internal/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func ScheduleRoutes(api fiber.Router) {
	schedule := api.Group("/schedule")
	//caching
	scheduleCache := cache.New(cache.Config{
		Expiration:   3 * time.Minute,
		CacheControl: true,
	})

	// Endpoint Publik (Bisa diakses user biasa tanpa login)
	schedule.Get("/", scheduleCache, handlers.GetAllSchedules)
	schedule.Get("/:id", scheduleCache, handlers.GetScheduleByID)

	// Endpoint Khusus Admin (Wajib login dan role = admin)
	schedule.Post("/", middleware.JWTProtected, middleware.AdminOnly, handlers.CreateSchedule)
	schedule.Patch("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.UpdateSchedule)
	schedule.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteSchedule)
}
