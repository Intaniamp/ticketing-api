package handlers

import (
	"database/sql"
	"ticketing-api/config"
	"ticketing-api/internal/models"
	"ticketing-api/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ProcessPayment godoc
//
//	@Summary        Proses Pembayaran Tiket
//	@Description    Melakukan simulasi pembayaran. Akan mengubah status booking menjadi 'success'.
//	@Tags           Payment
//	@Security       BearerAuth
//	@Accept         json
//	@Produce        json
//	@Param          request body        models.PaymentRequest   true    "Data Pembayaran"
//	@Success        201     {object}    models.Payment
//	@Failure        400     {object}    models.ErrorResponse
//	@Failure        403     {object}    models.ErrorResponse
//	@Failure        404     {object}    models.ErrorResponse
//	@Failure        500     {object}    models.ErrorResponse
//	@Router         /payment [post]
func ProcessPayment(c *fiber.Ctx) error {
	var req models.PaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	// 1. Ambil User ID dari JWT Token
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

	// 3. Validasi Keamanan
	if dbUserID != userID {
		tx.Rollback()
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{Message: "Akses ditolak: Ini bukan tiket Anda"})
	}
	if bookingStatus == "success" {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Tiket ini sudah dibayar sebelumnya"})
	}
	if bookingStatus == "failed" || bookingStatus == "expired" {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Waktu pembayaran telah habis, silakan pesan ulang"})
	}

	// 4. Masukkan data ke tabel payment
	newPaymentID := uuid.New().String()
	paymentDate := time.Now().Format("2006-01-02 15:04:05")
	paymentStatus := "success" // Simulasi langsung sukses

	insertQuery := "INSERT INTO payment (id, booking_id, amount, payment_method, payment_date, payment_status) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = tx.Exec(insertQuery, newPaymentID, req.BookingID, totalPrice, req.PaymentMethod, paymentDate, paymentStatus)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: "Gagal memproses pembayaran: " + err.Error()})
	}

	// 5. Update status di tabel booking menjadi success
	updateBookingQuery := "UPDATE booking SET status = 'success' WHERE id = ?"
	_, err = tx.Exec(updateBookingQuery, req.BookingID)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: "Gagal mengupdate status tiket: " + err.Error()})
	}

	// Commit semua perubahan jika berhasil
	if err = tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	// 🟢 6. TRIGGER WEBSOCKET DI SINI! (Setelah commit database sukses)
	utils.SendPaymentNotification(userID, "Hore! Pembayaran tiket kamu berhasil. Silakan cek detailnya di menu Tiket Saya.", "success")

	// 7. Kembalikan response sukses
	return c.Status(fiber.StatusCreated).JSON(models.Payment{
		ID:            newPaymentID,
		BookingID:     req.BookingID,
		Amount:        totalPrice,
		PaymentMethod: req.PaymentMethod,
		PaymentDate:   paymentDate,
		PaymentStatus: paymentStatus,
	})
}

// GetMyPayments godoc
//
//	@Summary		Ambil riwayat pembayaran user
//	@Description	Mengembalikan daftar semua pembayaran sukses/gagal milik user yang sedang login
//	@Tags			Payment
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.Payment
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/payment/my-history [get]
func GetMyPayments(c *fiber.Ctx) error {
	// 1. Ambil koper claims dari JWT
	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Message: "Unauthorized: Data token tidak valid"})
	}

	// 2. Buka koper claims, ambil "user_id"-nya
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Message: "Unauthorized: User ID tidak ditemukan dalam token"})
	}

	db := config.ConnectDB()
	defer db.Close()

	// 3. JOIN tabel payment dan booking untuk mengambil pembayaran milik user ini saja
	query := `
        SELECT p.id, p.booking_id, p.amount, p.payment_method, p.payment_date, p.payment_status 
        FROM payment p
        JOIN booking b ON p.booking_id = b.id
        WHERE b.user_id = ?
        ORDER BY p.payment_date DESC
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(&p.ID, &p.BookingID, &p.Amount, &p.PaymentMethod, &p.PaymentDate, &p.PaymentStatus); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
		}
		payments = append(payments, p)
	}

	if payments == nil {
		payments = []models.Payment{}
	}

	return c.Status(fiber.StatusOK).JSON(payments)
}

// GetAllPayments godoc
//
//	@Summary		Ambil semua data pembayaran (Khusus Admin)
//	@Description	Mengembalikan seluruh riwayat pembayaran dari semua user di dalam sistem
//	@Tags			Payment
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.Payment
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/payment [get]
func GetAllPayments(c *fiber.Ctx) error {
	db := config.ConnectDB()
	defer db.Close()

	query := "SELECT id, booking_id, amount, payment_method, payment_date, payment_status FROM payment ORDER BY payment_date DESC"
	rows, err := db.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(&p.ID, &p.BookingID, &p.Amount, &p.PaymentMethod, &p.PaymentDate, &p.PaymentStatus); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
		}
		payments = append(payments, p)
	}

	if payments == nil {
		payments = []models.Payment{}
	}

	return c.Status(fiber.StatusOK).JSON(payments)
}
