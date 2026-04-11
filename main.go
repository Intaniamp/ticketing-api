package main

import (
	"log"

	// Pastikan nama import ini sesuai dengan "go mod init ticketing-api"
	"ticketing-api/docs"
	"ticketing-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // Tambah CORS biar Flutter nggak error
	"github.com/gofiber/fiber/v2/middleware/logger"
	swagger "github.com/swaggo/fiber-swagger"
)

// @title			Ticketing API
// @version		1.0
// @description	API backend untuk aplikasi pemesanan tiket bioskop
// @host			localhost:3000
// @BasePath		/api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer <spasi> token-anda
func main() {
	// --- 1. Konfigurasi Swagger Programatik (INI YANG KAMU MAKSUD) ---
	// Ini gunanya biar info di UI Swagger muncul lengkap dan keren
	docs.SwaggerInfo.Title = "Ticketing API Cinema"
	docs.SwaggerInfo.Description = "API manajemen tiket bioskop - Project Multiplatform"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:3000"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// --- 2. Inisialisasi Fiber ---
	app := fiber.New(fiber.Config{
		AppName: "Ticketing API v1.0",
	})

	// --- 3. Middleware ---
	app.Use(logger.New()) // Buat ngeliat log di terminal
	app.Use(cors.New())   // PENTING: Supaya Flutter kamu nanti bisa akses API ini

	// --- 4. Route Swagger ---
	// Akses di: http://localhost:3000/swagger/index.html
	app.Get("/swagger/*", swagger.WrapHandler)

	// --- 5. Setup API Routes ---
	routes.SetupRoutes(app)

	// Redirect kalau buka localhost:3000 langsung ke Swagger
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	// --- 6. Jalankan Server ---
	log.Println("🚀 Server jalan di http://localhost:3000")
	log.Println("📖 Swagger UI: http://localhost:3000/swagger/index.html")
	log.Fatal(app.Listen(":3000"))
}
