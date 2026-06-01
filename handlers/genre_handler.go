package handlers

import (
	"database/sql"
	"ticketing-api/config"
	"ticketing-api/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllGenres godoc
//
//	@Summary		Ambil semua genre
//	@Description	Mengembalikan daftar seluruh genre film
//	@Tags			Genre
//	@Produce		json
//	@Success		200	{array}		models.Genre
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/genre [get]
func GetAllGenres(c *fiber.Ctx) error {
	db := config.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, genre_name FROM genre")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.GenreName); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
		}
		genres = append(genres, g)
	}

	return c.Status(fiber.StatusOK).JSON(genres)
}

// GetGenreByID godoc
//
//	@Summary		Ambil genre berdasarkan ID
//	@Description	Mengembalikan detail satu genre
//	@Tags			Genre
//	@Produce		json
//	@Param			id	path		int	true	"Genre ID"
//	@Success		200	{object}	models.Genre
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/genre/{id} [get]
func GetGenreByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := config.ConnectDB()
	defer db.Close()

	var g models.Genre
	err := db.QueryRow("SELECT id, genre_name FROM genre WHERE id = ?", id).Scan(&g.ID, &g.GenreName)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Genre tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(g)
}

// CreateGenre godoc
//
//	@Summary		Tambah genre baru
//	@Description	Menambahkan data genre baru ke dalam sistem
//	@Tags			Genre
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.GenreRequest	true	"Data Genre"
//	@Success		201		{object}	models.Genre
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/genre [post]
func CreateGenre(c *fiber.Ctx) error {
	var req models.GenreRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("INSERT INTO genre (genre_name) VALUES (?)", req.GenreName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	lastInsertID, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(models.Genre{
		ID:        int(lastInsertID),
		GenreName: req.GenreName,
	})
}

// UpdateGenre godoc
//
//	@Summary		Update genre
//	@Description	Mengubah nama genre berdasarkan ID
//	@Tags			Genre
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Genre ID"
//	@Param			request	body		models.GenreRequest	true	"Data Genre Baru"
//	@Success		200		{object}	models.Genre
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Router			/genre/{id} [patch]
func UpdateGenre(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.GenreRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("UPDATE genre SET genre_name = ? WHERE id = ?", req.GenreName, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Genre tidak ditemukan"})
	}

	var updated models.Genre
	db.QueryRow("SELECT id, genre_name FROM genre WHERE id = ?", id).Scan(&updated.ID, &updated.GenreName)

	return c.Status(fiber.StatusOK).JSON(updated)
}

// DeleteGenre godoc
//
//	@Summary		Hapus genre
//	@Description	Menghapus genre berdasarkan ID
//	@Tags			Genre
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Genre ID"
//	@Success		204	{string}	string
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/genre/{id} [delete]
func DeleteGenre(c *fiber.Ctx) error {
	id := c.Params("id")
	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("DELETE FROM genre WHERE id = ?", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Genre tidak ditemukan"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
