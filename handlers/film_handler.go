package handlers

import (
	"database/sql"
	"ticketing-api/config"
	"ticketing-api/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllFilm godoc
//
//	@Summary		Ambil semua data film
//	@Description	Mengembalikan seluruh daftar film dari database
//	@Tags			Film
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.Film
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/film [get]
func GetAllFilm(c *fiber.Ctx) error {
	db := config.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, title, duration, rating FROM film")
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	var films []models.Film
	for rows.Next() {
		var f models.Film
		rows.Scan(&f.ID, &f.Title, &f.Duration, &f.Rating, &f.Synopsis)
		films = append(films, f)
	}

	return c.Status(200).JSON(films)
}

// GetFilmByID godoc
//
//	@Summary		Ambil data film berdasarkan ID
//	@Description	Mengembalikan detail satu film berdasarkan parameter ID
//	@Tags			Film
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Film ID"
//	@Success		200	{object}	models.Film
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/film/{id} [get]
func GetFilmByID(c *fiber.Ctx) error {
	id := c.Params("id")

	db := config.ConnectDB()
	defer db.Close()

	var film models.Film
	err := db.QueryRow("SELECT id, title, duration, rating, synopsis FROM film WHERE id = ?", id).Scan(
		&film.ID,
		&film.Title,
		&film.Duration,
		&film.Rating,
		&film.Synopsis,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ErrorResponse{Message: "Film tidak ditemukan"})
		}

		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(200).JSON(film)
}

// CreateFilm godoc
//
//	@Summary		Buat film baru
//	@Description	Menambahkan data film baru (khusus admin)
//	@Tags			Film
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.Film	true	"Data film"
//	@Success		201		{object}	models.Film
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/film [post]
func CreateFilm(c *fiber.Ctx) error {
	var f models.Film
	if err := c.BodyParser(&f); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	_, err := db.Exec("INSERT INTO film (title, duration, rating, synopsis) VALUES (?, ?, ?, ?)", f.Title, f.Duration, f.Rating, f.Synopsis)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(201).JSON(f)
}

// UpdateFilm godoc
//
//	@Summary		Update film
//	@Description	Memperbarui data film berdasarkan ID (khusus admin)
//	@Tags			Film
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int			true	"Film ID"
//	@Param			request	body		models.Film	true	"Data film"
//	@Success		200		{object}	models.Film
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/film/{id} [patch]
func UpdateFilm(c *fiber.Ctx) error {
	id := c.Params("id")
	var f models.Film
	if err := c.BodyParser(&f); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	_, err := db.Exec("UPDATE film SET title=?, duration=?, rating=?, synopsis=? WHERE id=?", f.Title, f.Duration, f.Rating, f.Synopsis, id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(f)
}

// DeleteFilm godoc
//
//	@Summary		Hapus film
//	@Description	Menghapus data film berdasarkan ID (khusus admin)
//	@Tags			Film
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Film ID"
//	@Success		204	{string}	string
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/film/{id} [delete]
func DeleteFilm(c *fiber.Ctx) error {
	id := c.Params("id")

	db := config.ConnectDB()
	defer db.Close()

	_, err := db.Exec("DELETE FROM film WHERE id=?", id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.SendStatus(204)
}
