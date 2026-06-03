package handlers

import (
	"database/sql"
	"strconv"
	"ticketing-api/config"
	"ticketing-api/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllCinemas godoc
//
//	@Summary		Ambil semua bioskop
//	@Description	Mengembalikan daftar seluruh bioskop
//	@Tags			Cinema
//	@Produce		json
//	@Success		200	{array}		models.Cinema
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/cinema [get]
func GetAllCinemas(c *fiber.Ctx) error {
	db := config.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, cinema_name, city_id, address FROM cinema")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}
	defer rows.Close()

	var cinemas []models.Cinema
	for rows.Next() {
		var cinema models.Cinema
		if err := rows.Scan(&cinema.ID, &cinema.CinemaName, &cinema.CityID, &cinema.Address); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
		}
		cinemas = append(cinemas, cinema)
	}

	return c.Status(fiber.StatusOK).JSON(cinemas)
}

// GetCinemaByID godoc
//
//	@Summary		Ambil bioskop berdasarkan ID
//	@Description	Mengembalikan detail satu bioskop
//	@Tags			Cinema
//	@Produce		json
//	@Param			id	path		int	true	"Cinema ID"
//	@Success		200	{object}	models.Cinema
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/cinema/{id} [get]
func GetCinemaByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := config.ConnectDB()
	defer db.Close()

	var cinema models.Cinema
	err := db.QueryRow("SELECT id, cinema_name, city_id, address FROM cinema WHERE id = ?", id).Scan(&cinema.ID, &cinema.CinemaName, &cinema.CityID, &cinema.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Bioskop tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(cinema)
}

// CreateCinema godoc
//
//	@Summary		Tambah bioskop baru
//	@Description	Menambahkan data bioskop baru ke dalam sistem (khusus admin)
//	@Tags			Cinema
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CinemaRequest	true	"Data Bioskop"
//	@Success		201		{object}	models.Cinema
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/cinema [post]
func CreateCinema(c *fiber.Ctx) error {
	var req models.CinemaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("INSERT INTO cinema (cinema_name, city_id, address) VALUES (?, ?, ?)", req.CinemaName, req.CityID, req.Address)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	lastInsertID, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(models.Cinema{
		ID:         int(lastInsertID),
		CinemaName: req.CinemaName,
		CityID:     req.CityID,
		Address:    req.Address,
	})
}

// UpdateCinema godoc
//
//	@Summary		Update bioskop
//	@Description	Memperbarui data bioskop berdasarkan ID (khusus admin)
//	@Tags			Cinema
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Cinema ID"
//	@Param			request	body		models.CinemaRequest	true	"Data Bioskop"
//	@Success		200		{object}	models.Cinema
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/cinema/{id} [patch]
func UpdateCinema(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.CinemaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("UPDATE cinema SET cinema_name = ?, city_id = ?, address = ? WHERE id = ?", req.CinemaName, req.CityID, req.Address, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Bioskop tidak ditemukan"})
	}

	cinemaID, _ := strconv.Atoi(id)
	return c.Status(fiber.StatusOK).JSON(models.Cinema{
		ID:         cinemaID,
		CinemaName: req.CinemaName,
		CityID:     req.CityID,
		Address:    req.Address,
	})
}

// DeleteCinema godoc
//
//	@Summary		Hapus bioskop
//	@Description	Menghapus data bioskop berdasarkan ID (khusus admin)
//	@Tags			Cinema
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Cinema ID"
//	@Success		204	{string}	string
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/cinema/{id} [delete]
func DeleteCinema(c *fiber.Ctx) error {
	id := c.Params("id")

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("DELETE FROM cinema WHERE id = ?", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Bioskop tidak ditemukan"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
