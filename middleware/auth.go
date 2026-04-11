package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key") // ganti rahasia sesuai env variable

// Fungsi untuk generate token
func GenerateJWT(userID int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 1 hari
	})
	return token.SignedString(jwtSecret)
}

// Middleware Auth
func JWTProtected(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Missing token"})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
	}

	c.Locals("user", token.Claims.(jwt.MapClaims))
	return c.Next()
}

// Middleware untuk verifikasi role admin
func AdminOnly(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden: Admin only"})
	}
	return c.Next()
}
