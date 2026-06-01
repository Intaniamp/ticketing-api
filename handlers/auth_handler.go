package handlers

import (
	"database/sql"
	"ticketing-api/config"
	"ticketing-api/middleware"
	"ticketing-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Login godoc
//
//	@Summary		User Login
//	@Description	Autentikasi user menggunakan email dan password untuk mendapatkan JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.LoginRequest	true	"Kredensial Login"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/login [post]
func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request"})
	}

	db := config.ConnectDB()
	defer db.Close()

	var user models.User
	err := db.QueryRow("SELECT id, email, password, role FROM user WHERE email = ?", req.Email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Role,
	)
	if err == sql.ErrNoRows {
		return c.Status(401).JSON(fiber.Map{"message": "User not found"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Database error: " + err.Error()})
	}

	// contoh: tanpa hashing dulu (idealnya bcrypt)
	if user.Password != req.Password {
		return c.Status(401).JSON(fiber.Map{"message": "Wrong password"})
	}

	token, err := middleware.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Token generation failed"})
	}

	return c.JSON(fiber.Map{"token": token})
}

// Register godoc
//
//	@Summary		User Register
//	@Description	Mendaftarkan user baru ke sistem
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.RegisterRequest	true	"Data Registrasi"
//	@Success		201		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/register [post]
func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request"})
	}

	db := config.ConnectDB()
	defer db.Close()

	// 1. Generate UUID baru
	newID := uuid.New().String()

	// 2. Insert ke database (Tambahkan kolom id dan value newID)
	query := "INSERT INTO user (id, name, email, phone, password, role) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, newID, req.Name, req.Email, req.Phone, req.Password, "user")

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mendaftarkan user: " + err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Registrasi berhasil!"})
}
