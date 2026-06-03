package handlers

import (
	"database/sql"
	"strconv"
	"ticketing-api/config"
	"ticketing-api/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllStudios godoc
//
//	@Summary		Ambil semua studio
//	@Description	Mengembalikan daftar seluruh studio di semua bioskop
//	@Tags			Studio
//	@Produce		json
//	@Success		200	{array}		models.Studio
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/studio [get]
func GetAllStudios(c *fiber.Ctx) error {
	db := config.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, cinema_id, studio_name FROM studio")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}
	defer rows.Close()

	var studios []models.Studio
	for rows.Next() {
		var studio models.Studio
		if err := rows.Scan(&studio.ID, &studio.CinemaID, &studio.StudioName); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
		}
		studios = append(studios, studio)
	}

	return c.Status(fiber.StatusOK).JSON(studios)
}

// GetStudioByID godoc
//
//	@Summary		Ambil studio berdasarkan ID
//	@Description	Mengembalikan detail satu studio
//	@Tags			Studio
//	@Produce		json
//	@Param			id	path		int	true	"Studio ID"
//	@Success		200	{object}	models.Studio
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/studio/{id} [get]
func GetStudioByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := config.ConnectDB()
	defer db.Close()

	var studio models.Studio
	err := db.QueryRow("SELECT id, cinema_id, studio_name FROM studio WHERE id = ?", id).Scan(&studio.ID, &studio.CinemaID, &studio.StudioName)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Studio tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(studio)
}

// CreateStudio godoc
//
//	@Summary		Tambah studio baru
//	@Description	Menambahkan data studio baru ke dalam sistem (khusus admin)
//	@Tags			Studio
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.StudioRequest	true	"Data Studio"
//	@Success		201		{object}	models.Studio
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/studio [post]
func CreateStudio(c *fiber.Ctx) error {
	var req models.StudioRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("INSERT INTO studio (cinema_id, studio_name) VALUES (?, ?)", req.CinemaID, req.StudioName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	lastInsertID, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(models.Studio{
		ID:         int(lastInsertID),
		CinemaID:   req.CinemaID,
		StudioName: req.StudioName,
	})
}

// UpdateStudio godoc
//
//	@Summary		Update studio
//	@Description	Memperbarui data studio berdasarkan ID (khusus admin)
//	@Tags			Studio
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Studio ID"
//	@Param			request	body		models.StudioRequest	true	"Data Studio"
//	@Success		200		{object}	models.Studio
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/studio/{id} [patch]
func UpdateStudio(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.StudioRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("UPDATE studio SET cinema_id = ?, studio_name = ? WHERE id = ?", req.CinemaID, req.StudioName, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Studio tidak ditemukan"})
	}

	studioID, _ := strconv.Atoi(id)
	return c.Status(fiber.StatusOK).JSON(models.Studio{
		ID:         studioID,
		CinemaID:   req.CinemaID,
		StudioName: req.StudioName,
	})
}

// DeleteStudio godoc
//
//	@Summary		Hapus studio
//	@Description	Menghapus data studio berdasarkan ID (khusus admin)
//	@Tags			Studio
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Studio ID"
//	@Success		204	{string}	string
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/studio/{id} [delete]
func DeleteStudio(c *fiber.Ctx) error {
	id := c.Params("id")

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("DELETE FROM studio WHERE id = ?", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Studio tidak ditemukan"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
