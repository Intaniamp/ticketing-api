package handlers

import (
	"database/sql"
	"strconv"
	"ticketing-api/config"
	"ticketing-api/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllSchedules godoc
//
//	@Summary		Ambil semua jadwal tayang
//	@Description	Mengembalikan daftar seluruh jadwal tayang film
//	@Tags			Schedule
//	@Produce		json
//	@Success		200	{array}		models.Schedule
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/schedule [get]
func GetAllSchedules(c *fiber.Ctx) error {
	db := config.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, film_id, studio_id, date, time, price FROM schedule")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		if err := rows.Scan(&schedule.ID, &schedule.FilmID, &schedule.StudioID, &schedule.Date, &schedule.Time, &schedule.Price); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
		}
		schedules = append(schedules, schedule)
	}

	return c.Status(fiber.StatusOK).JSON(schedules)
}

// GetScheduleByID godoc
//
//	@Summary		Ambil jadwal berdasarkan ID
//	@Description	Mengembalikan detail satu jadwal tayang
//	@Tags			Schedule
//	@Produce		json
//	@Param			id	path		int	true	"Schedule ID"
//	@Success		200	{object}	models.Schedule
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/schedule/{id} [get]
func GetScheduleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := config.ConnectDB()
	defer db.Close()

	var schedule models.Schedule
	err := db.QueryRow("SELECT id, film_id, studio_id, date, time, price FROM schedule WHERE id = ?", id).Scan(&schedule.ID, &schedule.FilmID, &schedule.StudioID, &schedule.Date, &schedule.Time, &schedule.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Jadwal tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(schedule)
}

// CreateSchedule godoc
//
//	@Summary		Buat jadwal baru
//	@Description	Menambahkan data jadwal tayang baru (khusus admin)
//	@Tags			Schedule
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.ScheduleRequest	true	"Data Jadwal"
//	@Success		201		{object}	models.Schedule
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/schedule [post]
func CreateSchedule(c *fiber.Ctx) error {
	var req models.ScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("INSERT INTO schedule (film_id, studio_id, date, time, price) VALUES (?, ?, ?, ?, ?)", req.FilmID, req.StudioID, req.Date, req.Time, req.Price)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	lastInsertID, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(models.Schedule{
		ID:       int(lastInsertID),
		FilmID:   req.FilmID,
		StudioID: req.StudioID,
		Date:     req.Date,
		Time:     req.Time,
		Price:    req.Price,
	})
}

// UpdateSchedule godoc
//
//	@Summary		Update jadwal
//	@Description	Memperbarui data jadwal tayang berdasarkan ID (khusus admin)
//	@Tags			Schedule
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Schedule ID"
//	@Param			request	body		models.ScheduleRequest	true	"Data Jadwal"
//	@Success		200		{object}	models.Schedule
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/schedule/{id} [patch]
func UpdateSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.ScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("UPDATE schedule SET film_id = ?, studio_id = ?, date = ?, time = ?, price = ? WHERE id = ?", req.FilmID, req.StudioID, req.Date, req.Time, req.Price, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Jadwal tidak ditemukan"})
	}

	scheduleID, _ := strconv.Atoi(id)
	return c.Status(fiber.StatusOK).JSON(models.Schedule{
		ID:       scheduleID,
		FilmID:   req.FilmID,
		StudioID: req.StudioID,
		Date:     req.Date,
		Time:     req.Time,
		Price:    req.Price,
	})
}

// DeleteSchedule godoc
//
//	@Summary		Hapus jadwal
//	@Description	Menghapus data jadwal tayang berdasarkan ID (khusus admin)
//	@Tags			Schedule
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Schedule ID"
//	@Success		204	{string}	string
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/schedule/{id} [delete]
func DeleteSchedule(c *fiber.Ctx) error {
	id := c.Params("id")

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("DELETE FROM schedule WHERE id = ?", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Jadwal tidak ditemukan"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
