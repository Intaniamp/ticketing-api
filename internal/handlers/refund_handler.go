package handlers

import (
	"database/sql"
	"ticketing-api/config"
	"ticketing-api/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ProcessRefund godoc
//
//	@Summary		Proses Pembatalan dan Refund Tiket
//	@Description	Membatalkan tiket, mencatat refund, dan mengosongkan kembali kursi yang sudah dipesan.
//	@Tags			Refund
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.RefundRequest	true	"Data Refund"
//	@Success		201		{object}	models.Refund
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		403		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/refund [post]
func ProcessRefund(c *fiber.Ctx) error {
	var req models.RefundRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	// 1. Validasi Token JWT
	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Message: "Unauthorized: Data token tidak valid"})
	}
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Message: "Unauthorized: User ID tidak ditemukan"})
	}

	db := config.ConnectDB()
	defer db.Close()

	// Memulai Database Transaction
	tx, err := db.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	// 2. Cek status booking dan harganya (pastikan milik user yang sedang login)
	var bookingStatus string
	var totalPrice float64
	var dbUserID string

	err = tx.QueryRow("SELECT status, total_price, user_id FROM booking WHERE id = ?", req.BookingID).Scan(&bookingStatus, &totalPrice, &dbUserID)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Data booking tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	// 3. Validasi Keamanan Bisnis
	if dbUserID != userID {
		tx.Rollback()
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Message: "Akses ditolak: Anda tidak bisa merefund tiket orang lain"})
	}
	if bookingStatus != "success" {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Hanya tiket yang sudah dibayar (sukses) yang bisa direfund"})
	}

	// 4. Masukkan data ke tabel refund
	newRefundID := uuid.New().String()
	refundDate := time.Now().Format("2006-01-02 15:04:05")

	insertRefundQuery := "INSERT INTO refund (id, booking_id, refund_date, refund_amount) VALUES (?, ?, ?, ?)"
	_, err = tx.Exec(insertRefundQuery, newRefundID, req.BookingID, refundDate, totalPrice)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: "Gagal memproses pencatatan refund: " + err.Error()})
	}

	// 5. Update status di tabel booking menjadi 'refunded'
	updateBookingQuery := "UPDATE booking SET status = 'refunded' WHERE id = ?"
	_, err = tx.Exec(updateBookingQuery, req.BookingID)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: "Gagal mengupdate status tiket: " + err.Error()})
	}

	// 6. LEPASKAN KURSI (Hapus data dari tabel pivot booking_seat)
	// Inilah keajaibannya, dengan dihapus, kursi ini otomatis kembali 'Available'
	deleteSeatQuery := "DELETE FROM booking_seat WHERE booking_id = ?"
	_, err = tx.Exec(deleteSeatQuery, req.BookingID)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: "Gagal melepaskan kursi tiket: " + err.Error()})
	}

	// Commit semua perubahan jika 3 proses di atas sukses
	if err = tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	// 7. Kembalikan response
	return c.Status(fiber.StatusCreated).JSON(models.Refund{
		ID:           newRefundID,
		BookingID:    req.BookingID,
		RefundDate:   refundDate,
		RefundAmount: totalPrice,
	})
}
