package handlers

import (
	"database/sql"
	"strings"
	"ticketing-api/config"
	"ticketing-api/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const bookingDateTimeLayout = "2006-01-02 15:04:05"

func fetchBookings(db *sql.DB, query string, args ...any) ([]models.Booking, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]models.Booking, 0)
	for rows.Next() {
		var booking models.Booking
		if err := rows.Scan(&booking.ID, &booking.UserID, &booking.ScheduleID, &booking.BookingDate, &booking.ExpiredAt, &booking.Status, &booking.TotalPrice); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, rows.Err()
}

// CreateBooking godoc
//
//	@Summary		Buat pesanan tiket (Booking)
//	@Description	Memesan kursi. Status awal otomatis 'pending'. Menerapkan penguncian kursi (Seat Locking).
//	@Tags			Booking
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.BookingRequest	true	"Data Booking"
//	@Success		201		{object}	models.Booking
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		409		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/booking [post]
func CreateBooking(c *fiber.Ctx) error {
	var req models.BookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	if len(req.SeatIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Minimal harus memilih 1 kursi"})
	}

	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Message: "Unauthorized: Data token tidak valid"})
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Message: "Unauthorized: User ID tidak ditemukan dalam token"})
	}

	db := config.ConnectDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	var ticketPrice float64
	err = tx.QueryRow("SELECT price FROM schedule WHERE id = ?", req.ScheduleID).Scan(&ticketPrice)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Jadwal tayang tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	for _, seatID := range req.SeatIDs {
		var count int
		checkQuery := `
			SELECT COUNT(*) 
			FROM booking_seat bs
			JOIN booking b ON bs.booking_id = b.id
			WHERE b.schedule_id = ? AND bs.seat_id = ? AND b.status IN ('pending', 'success')
		`
		err = tx.QueryRow(checkQuery, req.ScheduleID, seatID).Scan(&count)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
		}

		if count > 0 {
			tx.Rollback()
			return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{Message: "Mohon maaf, kursi pilihanmu sudah dipesan orang lain"})
		}
	}

	newBookingID := uuid.New().String()
	now := time.Now()
	expiredTime := now.Add(5 * time.Minute)
	bookingDate := now.Format(bookingDateTimeLayout)
	expiredAt := expiredTime.Format(bookingDateTimeLayout)

	totalPrice := ticketPrice * float64(len(req.SeatIDs))
	initialStatus := "pending"

	insertBookingQuery := "INSERT INTO booking (id, user_id, schedule_id, booking_date, expired_at, status, total_price) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err = tx.Exec(insertBookingQuery, newBookingID, userID, req.ScheduleID, bookingDate, expiredAt, initialStatus, totalPrice)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	// 4. Masukkan data ke tabel pivot booking_seat
	for _, seatID := range req.SeatIDs {
		insertPivotQuery := "INSERT INTO booking_seat (booking_id, seat_id) VALUES (?, ?)"
		_, err = tx.Exec(insertPivotQuery, newBookingID, seatID)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: "Gagal menyimpan relasi kursi: " + err.Error()})
		}
	}

	if err = tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Booking{
		ID:          newBookingID,
		UserID:      userID,
		ScheduleID:  req.ScheduleID,
		BookingDate: bookingDate,
		Status:      initialStatus,
		TotalPrice:  totalPrice,
		ExpiredAt:   expiredAt,
	})
}

// GetMyBookings godoc
//
//	@Summary		Ambil riwayat pesanan user
//	@Description	Mengembalikan daftar pesanan ringkas milik user yang sedang login
//	@Tags			Booking
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.Booking
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/booking/my-history [get]
func GetMyBookings(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Message: "Unauthorized: Data token tidak valid"})
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Message: "Unauthorized: User ID tidak ditemukan dalam token"})
	}

	db := config.ConnectDB()
	defer db.Close()

	query := "SELECT id, user_id, schedule_id, booking_date, COALESCE(expired_at, ''), status, total_price FROM booking WHERE user_id = ? ORDER BY booking_date DESC"
	bookings, err := fetchBookings(db, query, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(bookings)
}

// GetAllBookings godoc
//
//	@Summary		Ambil semua pesanan (Khusus Admin)
//	@Description	Mengembalikan seluruh daftar booking dari semua user di dalam sistem
//	@Tags			Booking
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.Booking
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/booking [get]
func GetAllBookings(c *fiber.Ctx) error {
	db := config.ConnectDB()
	defer db.Close()

	query := "SELECT id, user_id, schedule_id, booking_date, COALESCE(expired_at, ''), status, total_price FROM booking ORDER BY booking_date DESC"
	bookings, err := fetchBookings(db, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(bookings)
}

// GetBookingByID godoc
//
//	@Summary		Ambil detail E-Ticket
//	@Description	Mengembalikan detail lengkap satu tiket (Film, Bioskop, Jadwal, Kursi)
//	@Tags			Booking
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"Booking ID (UUID)"
//	@Success		200	{object}	models.BookingDetailResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/booking/{id} [get]
func GetBookingByID(c *fiber.Ctx) error {
	bookingID := c.Params("id")

	db := config.ConnectDB()
	defer db.Close()

	var res models.BookingDetailResponse
	res.ID = bookingID

	uuidParts := strings.Split(bookingID, "-")
	if len(uuidParts) > 0 {
		res.OrderID = "#CM-" + strings.ToUpper(uuidParts[0])
	} else {
		res.OrderID = "#CM-" + bookingID
	}

	// 1. Ambil data gabungan dari 5 tabel
	query := `
		SELECT 
			b.status, b.total_price, COALESCE(b.expired_at, ''),
			f.title, COALESCE(f.poster, ''),
			c.cinema_name,
			sch.date, sch.time
		FROM booking b
		JOIN schedule sch ON b.schedule_id = sch.id
		JOIN film f ON sch.film_id = f.id
		JOIN studio st ON sch.studio_id = st.id
		JOIN cinema c ON st.cinema_id = c.id
		WHERE b.id = ?
	`

	err := db.QueryRow(query, bookingID).Scan(
		&res.Status, &res.TotalPrice, &res.ExpiredAt,
		&res.FilmTitle, &res.PosterURL,
		&res.CinemaName,
		&res.Date, &res.Time,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Data booking tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	// 2. Ambil daftar kursi
	seatQuery := `
		SELECT s.seat_number 
		FROM seat s
		JOIN booking_seat bs ON s.id = bs.seat_id
		WHERE bs.booking_id = ?
	`
	rows, err := db.Query(seatQuery, bookingID)
	res.Seats = []string{} // Inisialisasi awal agar tidak null di JSON

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var seatNum string
			if err := rows.Scan(&seatNum); err == nil {
				res.Seats = append(res.Seats, seatNum)
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
