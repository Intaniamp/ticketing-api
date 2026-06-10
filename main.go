package main

import (
	"log"

	"ticketing-api/docs"
	"ticketing-api/routes"
	"ticketing-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2" 
	"github.com/joho/godotenv"
	swagger "github.com/swaggo/fiber-swagger"
)

// @title          Ticketing API
// @version        1.0
// @description    API backend untuk aplikasi pemesanan tiket bioskop
// @host           localhost:3000
// @BasePath       /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer <spasi> token-anda
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// --- 1. Konfigurasi Swagger ---
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

	// 🟢 TAMBAHAN: Jalankan Hub WebSocket di background (Wajib sebelum rute lain)
	go utils.StartWebSocketHub()

	// 🟢 3. MIDDLEWARE HARUS DI SINI (Paling Atas)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH",
	}))

	// 🟢 TAMBAHAN: Middleware khusus mengecek koneksi WebSocket (Biar ga crash)
		// 🟢 REVISI: Tambahkan /:userID di akhir rute agar tahu ini milik siapa
	app.Get("/ws/notifications/:userID", websocket.New(func(c *websocket.Conn) {
		// 1. Ambil userID dari parameter URL
		userID := c.Params("userID")
		if userID == "" {
			log.Println("⚠️ Koneksi ditolak: User ID tidak ditemukan di URL")
			c.Close()
			return
		}

		// 2. Buat struktur objek client baru menggunakan struct yang sudah kita buat di utils
		newClient := &utils.Client{
			Conn:   c,
			UserID: userID,
		}

		// 3. Masukkan ke dalam map dengan aman menggunakan Mutex
		utils.WSManager.Mutex.Lock()
		utils.WSManager.Clients[newClient] = true
		utils.WSManager.Mutex.Unlock()

		log.Printf("🔌 Client dengan UserID %s berhasil terkoneksi ke WebSocket!\n", userID)

		defer func() {
			utils.WSManager.Mutex.Lock()
			delete(utils.WSManager.Clients, newClient)
			utils.WSManager.Mutex.Unlock()
			c.Close()
			log.Printf("❌ Client dengan UserID %s disconnect dari WebSocket.\n", userID)
		}()

		// 5. Keep-alive loop: Menjaga koneksi tetap hidup
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break // Keluar loop jika koneksi terputus (misal: HP mati atau pindah halaman)
			}
		}
	}))

	// 🟢 4. SETELAH CORS, BARU BUKA FOLDER STATIC (Dengan Cache Dimatikan!)
	app.Static("/uploads", "./public/uploads")
	app.Static("/", "./public")

	// --- 5. Route Swagger ---
	// Akses di: http://localhost:3000/swagger/index.html
	app.Get("/swagger/*", swagger.WrapHandler)

	// --- 6. Setup API Routes ---
	utils.StartTicketSweeper()

	// Setup API Routes
	routes.SetupRoutes(app)

	// Redirect kalau buka localhost:3000 langsung ke Swagger
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	// --- 7. Jalankan Server ---
	log.Println("🚀 Server jalan di http://localhost:3000")
	log.Println("📖 Swagger UI: http://localhost:3000/swagger/index.html")
	log.Fatal(app.Listen(":3000"))
}