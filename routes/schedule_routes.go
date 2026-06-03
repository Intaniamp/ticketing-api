package routes

import (
	"ticketing-api/handlers"
	"ticketing-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func ScheduleRoutes(api fiber.Router) {
	schedule := api.Group("/schedule")

	// Endpoint Publik (Bisa diakses user biasa tanpa login)
	schedule.Get("/", handlers.GetAllSchedules)
	schedule.Get("/:id", handlers.GetScheduleByID)

	// Endpoint Khusus Admin (Wajib login dan role = admin)
	schedule.Post("/", middleware.JWTProtected, middleware.AdminOnly, handlers.CreateSchedule)
	schedule.Patch("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.UpdateSchedule)
	schedule.Delete("/:id", middleware.JWTProtected, middleware.AdminOnly, handlers.DeleteSchedule)
}
