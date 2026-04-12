package handlers

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"ticketing-api/config"
	"ticketing-api/models"
	"time"

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

	// Ambil semua kolom yang ada di struct Film
	rows, err := db.Query("SELECT id, title, duration, synopsis, poster_url, age_rating, release_year FROM film")
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	var films []models.Film
	for rows.Next() {
		var f models.Film
		// Scan semua kolom secara berurutan
		rows.Scan(&f.ID, &f.Title, &f.Duration, &f.Synopsis, &f.PosterURL, &f.AgeRating, &f.ReleaseYear)
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
//
// @Param id path int true "Film ID"
//
//	@Success		200	{object}	models.Film
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/film/{id} [get]
func GetFilmByID(c *fiber.Ctx) error {
	id := c.Params("id")

	db := config.ConnectDB()
	defer db.Close()

	var film models.Film
	err := db.QueryRow("SELECT id, title, duration, synopsis, poster_url, age_rating, release_year FROM film WHERE id = ?", id).Scan(
		&film.ID,
		&film.Title,
		&film.Duration,
		&film.Synopsis,
		&film.PosterURL,
		&film.AgeRating,
		&film.ReleaseYear,
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
//
// @Param request body models.FilmRequest true "Data film"
//
//	@Success		201		{object}	models.Film
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/film [post]
func CreateFilm(c *fiber.Ctx) error {
	var req models.FilmRequest // Pakai struct Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	// Update query sesuai kolom baru di DB
	query := "INSERT INTO film (title, duration, synopsis, age_rating, release_year) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, req.Title, req.Duration, req.Synopsis, req.AgeRating, req.ReleaseYear)

	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(201).JSON(req)
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
//	@Param			request	body		models.FilmRequest	true	"Data film"
//	@Success		200		{object}	models.Film
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/film/{id} [patch]
func UpdateFilm(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.FilmRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	query := "UPDATE film SET title=?, duration=?, synopsis=?, age_rating=?, release_year=? WHERE id=?"
	_, err := db.Exec(query, req.Title, req.Duration, req.Synopsis, req.AgeRating, req.ReleaseYear, id)

	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(req)
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

// UploadPoster godoc
//
//  @Summary        Upload Poster Film
//  @Description    Mengupload file poster untuk film tertentu (khusus admin)
//  @Tags           Film
//  @Security       BearerAuth
//  @Accept         multipart/form-data
//  @Produce        json
// @Param id path int true "Film ID"
// @Param poster formData file true "File Gambar"
//  @Success 200 {object} map[string]string
//  @Failure 400 {object} models.ErrorResponse
//  @Failure 500 {object} models.ErrorResponse
//  @Router /film/{id}/poster [post]
func UploadPoster(c *fiber.Ctx) error {
	id := c.Params("id")

	file, err := c.FormFile("poster")
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Message: "File poster tidak ditemukan"})
	}

	var maxFileSize int64 = 2 * 1024 * 1024 // 2MB
	if file.Size > maxFileSize {
		return c.Status(400).JSON(models.ErrorResponse{Message: "Ukuran file terlalu besar! Maksimal 2MB"})
	}

	extension := filepath.Ext(file.Filename)
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowedExtensions[extension] {
		return c.Status(400).JSON(models.ErrorResponse{Message: "Format file tidak didukung! Gunakan jpg, jpeg, atau png"})
	}

	newFileName := fmt.Sprintf("%d%s", time.Now().Unix(), extension)

	savePath := filepath.Join("public/uploads/posters", newFileName)
	dbPath := fmt.Sprintf("uploads/posters/%s", newFileName)

	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: "Gagal menyimpan file di server"})
	}

	db := config.ConnectDB()
	defer db.Close()

	_, err = db.Exec("UPDATE film SET poster_url = ? WHERE id = ?", dbPath, id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":    "Poster berhasil diupload",
		"poster_url": dbPath,
	})
}
