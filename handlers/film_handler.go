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
//	@Description	Mengembalikan seluruh daftar film beserta genrenya
//	@Tags			Film
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.Film
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/film [get]
func GetAllFilm(c *fiber.Ctx) error {
	db := config.ConnectDB()
	defer db.Close()

	query := "SELECT id, title, duration, COALESCE(synopsis, ''), COALESCE(poster, ''), COALESCE(age_rating, ''), release_year FROM film"
	rows, err := db.Query(query)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}
	defer rows.Close()

	var films []models.Film
	for rows.Next() {
		var f models.Film
		if err := rows.Scan(&f.ID, &f.Title, &f.Duration, &f.Synopsis, &f.PosterURL, &f.AgeRating, &f.ReleaseYear); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
		}

		// Ambil genre untuk masing-masing film
		genreQuery := `
			SELECT g.id, g.genre_name 
			FROM genre g
			JOIN genre_film gf ON g.id = gf.genre_id
			WHERE gf.film_id = ?`

		genreRows, _ := db.Query(genreQuery, f.ID)
		f.Genres = []models.Genre{} // Inisialisasi array kosong agar tidak null di JSON

		for genreRows.Next() {
			var g models.Genre
			if err := genreRows.Scan(&g.ID, &g.GenreName); err == nil {
				f.Genres = append(f.Genres, g)
			}
		}
		genreRows.Close()

		films = append(films, f)
	}

	return c.Status(200).JSON(films)
}

// GetFilmByID godoc
//
//	@Summary		Ambil data film berdasarkan ID
//	@Description	Mengembalikan detail satu film beserta genrenya
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
	query := "SELECT id, title, duration, COALESCE(synopsis, ''), COALESCE(poster, ''), COALESCE(age_rating, ''), release_year FROM film WHERE id = ?"
	err := db.QueryRow(query, id).Scan(
		&film.ID, &film.Title, &film.Duration, &film.Synopsis, &film.PosterURL, &film.AgeRating, &film.ReleaseYear,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ErrorResponse{Message: "Film tidak ditemukan"})
		}
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	// Ambil list genre
	genreQuery := `
		SELECT g.id, g.genre_name 
		FROM genre g
		JOIN genre_film gf ON g.id = gf.genre_id
		WHERE gf.film_id = ?`

	rows, err := db.Query(genreQuery, film.ID)
	film.Genres = []models.Genre{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var g models.Genre
			if err := rows.Scan(&g.ID, &g.GenreName); err == nil {
				film.Genres = append(film.Genres, g)
			}
		}
	}

	return c.Status(200).JSON(film)
}

// CreateFilm godoc
//
//	@Summary		Buat film baru
//	@Description	Menambahkan data film baru beserta genrenya (khusus admin)
//	@Tags			Film
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.FilmRequest	true	"Data film"
//	@Success		201		{object}	models.Film
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/film [post]
func CreateFilm(c *fiber.Ctx) error {
	var req models.FilmRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	// 1. Insert tabel film
	query := "INSERT INTO film (title, duration, synopsis, age_rating, release_year) VALUES (?, ?, ?, ?, ?)"
	result, err := db.Exec(query, req.Title, req.Duration, req.Synopsis, req.AgeRating, req.ReleaseYear)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	lastInsertID, _ := result.LastInsertId()
	filmID := int(lastInsertID)

	// 2. Insert tabel pivot genre_film
	var assignedGenres []models.Genre
	for _, genreID := range req.GenreIDs {
		_, err := db.Exec("INSERT INTO genre_film (film_id, genre_id) VALUES (?, ?)", filmID, genreID)
		if err == nil {
			var g models.Genre
			db.QueryRow("SELECT id, genre_name FROM genre WHERE id = ?", genreID).Scan(&g.ID, &g.GenreName)
			assignedGenres = append(assignedGenres, g)
		}
	}

	if assignedGenres == nil {
		assignedGenres = []models.Genre{}
	}

	return c.Status(201).JSON(models.Film{
		ID:          filmID,
		Title:       req.Title,
		Duration:    req.Duration,
		Synopsis:    req.Synopsis,
		AgeRating:   req.AgeRating,
		ReleaseYear: req.ReleaseYear,
		PosterURL:   "",
		Genres:      assignedGenres,
	})
}

// UpdateFilm godoc
//
//	@Summary		Update film
//	@Description	Memperbarui data film dan relasi genrenya (khusus admin)
//	@Tags			Film
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Film ID"
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

	// 1. Update tabel film
	query := "UPDATE film SET title=?, duration=?, synopsis=?, age_rating=?, release_year=? WHERE id=?"
	result, err := db.Exec(query, req.Title, req.Duration, req.Synopsis, req.AgeRating, req.ReleaseYear, id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(404).JSON(models.ErrorResponse{Message: "Film tidak ditemukan"})
	}

	// 2. Hapus genre lama, masukkan genre baru
	db.Exec("DELETE FROM genre_film WHERE film_id = ?", id)
	for _, genreID := range req.GenreIDs {
		db.Exec("INSERT INTO genre_film (film_id, genre_id) VALUES (?, ?)", id, genreID)
	}

	// 3. Ambil data terbaru untuk dikembalikan ke FE
	var updated models.Film
	selectQuery := "SELECT id, title, duration, COALESCE(synopsis, ''), COALESCE(poster, ''), COALESCE(age_rating, ''), release_year FROM film WHERE id = ?"
	db.QueryRow(selectQuery, id).Scan(
		&updated.ID, &updated.Title, &updated.Duration, &updated.Synopsis, &updated.PosterURL, &updated.AgeRating, &updated.ReleaseYear,
	)

	// Ambil list genre terbaru
	genreQuery := `
		SELECT g.id, g.genre_name 
		FROM genre g
		JOIN genre_film gf ON g.id = gf.genre_id
		WHERE gf.film_id = ?`

	rows, _ := db.Query(genreQuery, id)
	updated.Genres = []models.Genre{}
	defer rows.Close()
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.GenreName); err == nil {
			updated.Genres = append(updated.Genres, g)
		}
	}

	return c.Status(200).JSON(updated)
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

	// Relasi di genre_film akan otomatis terhapus karena ON DELETE CASCADE di database
	result, err := db.Exec("DELETE FROM film WHERE id=?", id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(404).JSON(models.ErrorResponse{Message: "Film tidak ditemukan"})
	}

	return c.SendStatus(204)
}

// UploadPoster godoc
//
//	@Summary		Upload Poster Film
//	@Description	Mengupload file poster untuk film tertentu (khusus admin)
//	@Tags			Film
//	@Security		BearerAuth
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		int		true	"Film ID"
//	@Param			poster	formData	file	true	"File Gambar"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/film/{id}/poster [post]
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

	_, err = db.Exec("UPDATE film SET poster = ? WHERE id = ?", dbPath, id)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":    "Poster berhasil diupload",
		"poster_url": dbPath,
	})
}
