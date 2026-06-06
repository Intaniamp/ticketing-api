package handlers

import (
	"ticketing-api/config"
	"ticketing-api/internal/models"

	"github.com/gofiber/fiber/v2"
)

// GetSeatsByStudio godoc
//
//	@Summary		Ambil kursi berdasarkan Studio
//	@Description	Mengembalikan daftar seluruh kursi yang ada di dalam studio tertentu
//	@Tags			Seat
//	@Produce		json
//	@Param			studio_id	path		int	true	"Studio ID"
//	@Success		200			{array}		models.Seat
//	@Failure		500			{object}	models.ErrorResponse
//	@Router			/seat/studio/{studio_id} [get]
func GetSeatsByStudio(c *fiber.Ctx) error {
	studioID := c.Params("studio_id")
	db := config.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, studio_id, seat_number FROM seat WHERE studio_id = ?", studioID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var s models.Seat
		if err := rows.Scan(&s.ID, &s.StudioID, &s.SeatNumber); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
		}
		seats = append(seats, s)
	}

	// Kalau kosong, kembalikan array kosong, bukan null
	if seats == nil {
		seats = []models.Seat{}
	}

	return c.Status(fiber.StatusOK).JSON(seats)
}

// BulkCreateSeats godoc
//
//	@Summary		Tambah banyak kursi sekaligus (Bulk)
//	@Description	Menambahkan banyak kursi ke dalam satu studio dalam satu kali request (khusus admin)
//	@Tags			Seat
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.BulkSeatRequest	true	"Data Kursi Massal"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/seat/bulk [post]
func BulkCreateSeats(c *fiber.Ctx) error {
	var req models.BulkSeatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	if len(req.SeatNumbers) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Minimal harus ada 1 nomor kursi"})
	}

	db := config.ConnectDB()
	defer db.Close()

	// Memulai Database Transaction
	tx, err := db.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	// Looping insert untuk setiap nomor kursi
	for _, seatNo := range req.SeatNumbers {
		_, err := tx.Exec("INSERT INTO seat (studio_id, seat_number) VALUES (?, ?)", req.StudioID, seatNo)
		if err != nil {
			// Kalau error (misal duplicate atau studio gak ada), batalkan semua!
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: "Gagal menyimpan data: " + err.Error()})
		}
	}

	// Commit transaksi jika semua berhasil
	err = tx.Commit()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":           "Berhasil menambahkan kursi secara massal",
		"studio_id":         req.StudioID,
		"total_seats_added": len(req.SeatNumbers),
	})
}

// DeleteSeat godoc
//
//	@Summary		Hapus kursi
//	@Description	Menghapus kursi berdasarkan ID (khusus admin)
//	@Tags			Seat
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Seat ID"
//	@Success		204	{string}	string
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/seat/{id} [delete]
func DeleteSeat(c *fiber.Ctx) error {
	id := c.Params("id")
	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("DELETE FROM seat WHERE id = ?", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Kursi tidak ditemukan"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// DeleteSeatsByStudio godoc
//
//	@Summary		Sapu bersih kursi di satu studio
//	@Description	Menghapus seluruh kursi berdasarkan Studio ID (khusus admin)
//	@Tags			Seat
//	@Security		BearerAuth
//	@Param			studio_id	path		int	true	"Studio ID"
//	@Success		200			{object}	map[string]interface{}
//	@Failure		500			{object}	models.ErrorResponse
//	@Router			/seat/studio/{studio_id} [delete]
func DeleteSeatsByStudio(c *fiber.Ctx) error {
	studioID := c.Params("studio_id")
	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("DELETE FROM seat WHERE studio_id = ?", studioID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Tidak ada kursi yang ditemukan di studio ini"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Seluruh kursi di studio berhasil dihapus",
		"studio_id":     studioID,
		"deleted_count": affected,
	})
}

// GetSeatsBySchedule godoc
//
//	@Summary		Lihat Denah Kursi Live
//	@Description	Mengembalikan daftar semua kursi beserta statusnya (available/booked) untuk jadwal tayang tertentu. Ini yang dipakai User!
//	@Tags			Seat
//	@Produce		json
//	@Param			schedule_id	path		int	true	"ID Jadwal Tayang"
//	@Success		200			{array}		models.SeatStatus
//	@Failure		500			{object}	models.ErrorResponse
//	@Router			/seat/schedule/{schedule_id} [get]
func GetSeatsBySchedule(c *fiber.Ctx) error {
	scheduleID := c.Params("schedule_id")

	db := config.ConnectDB()
	defer db.Close()

	// Query ajaib untuk langsung memisahkan kursi yang kosong dan yang sudah dipesan
	query := `
		SELECT s.id, s.seat_number,
			CASE 
				WHEN s.id IN (
					SELECT bs.seat_id 
					FROM booking_seat bs
					JOIN booking b ON bs.booking_id = b.id
					WHERE b.schedule_id = ? AND b.status IN ('pending', 'success')
				) THEN 'booked' 
				ELSE 'available' 
			END as status
		FROM seat s
		WHERE s.studio_id = (SELECT studio_id FROM schedule WHERE id = ?)
		ORDER BY s.id ASC
	`

	rows, err := db.Query(query, scheduleID, scheduleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}
	defer rows.Close()

	var seats []models.SeatStatus
	for rows.Next() {
		var seat models.SeatStatus
		if err := rows.Scan(&seat.ID, &seat.SeatNumber, &seat.Status); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
		}
		seats = append(seats, seat)
	}

	if seats == nil {
		seats = []models.SeatStatus{}
	}

	return c.Status(fiber.StatusOK).JSON(seats)
}
