package handlers

import (
	"database/sql"
	"strconv"
	"ticketing-api/config"
	"ticketing-api/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllCities godoc
//
//	@Summary		Ambil semua kota
//	@Description	Mengembalikan daftar seluruh kota
//	@Tags			City
//	@Produce		json
//	@Success		200	{array}		models.City
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/city [get]
func GetAllCities(c *fiber.Ctx) error {
	db := config.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, city_name FROM city")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}
	defer rows.Close()

	var cities []models.City
	for rows.Next() {
		var city models.City
		if err := rows.Scan(&city.ID, &city.CityName); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
		}
		cities = append(cities, city)
	}

	return c.Status(fiber.StatusOK).JSON(cities)
}

// GetCityByID godoc
//
//	@Summary		Ambil kota berdasarkan ID
//	@Description	Mengembalikan detail satu kota
//	@Tags			City
//	@Produce		json
//	@Param			id	path		int	true	"City ID"
//	@Success		200	{object}	models.City
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/city/{id} [get]
func GetCityByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := config.ConnectDB()
	defer db.Close()

	var city models.City
	err := db.QueryRow("SELECT id, city_name FROM city WHERE id = ?", id).Scan(&city.ID, &city.CityName)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Kota tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(city)
}

// CreateCity godoc
//
//	@Summary		Tambah kota baru
//	@Description	Menambahkan data kota baru ke dalam sistem (khusus admin)
//	@Tags			City
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CityRequest	true	"Data Kota"
//	@Success		201		{object}	models.City
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/city [post]
func CreateCity(c *fiber.Ctx) error {
	var req models.CityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("INSERT INTO city (city_name) VALUES (?)", req.CityName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	lastInsertID, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(models.City{
		ID:       int(lastInsertID),
		CityName: req.CityName,
	})
}

// UpdateCity godoc
//
//	@Summary		Update kota
//	@Description	Memperbarui data kota berdasarkan ID (khusus admin)
//	@Tags			City
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"City ID"
//	@Param			request	body		models.CityRequest	true	"Data Kota"
//	@Success		200		{object}	models.City
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/city/{id} [patch]
func UpdateCity(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.CityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("UPDATE city SET city_name = ? WHERE id = ?", req.CityName, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Kota tidak ditemukan"})
	}

	cityID, _ := strconv.Atoi(id)
	return c.Status(fiber.StatusOK).JSON(models.City{
		ID:       cityID,
		CityName: req.CityName,
	})
}

// DeleteCity godoc
//
//	@Summary		Hapus kota
//	@Description	Menghapus data kota berdasarkan ID (khusus admin)
//	@Tags			City
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"City ID"
//	@Success		204	{string}	string
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/city/{id} [delete]
func DeleteCity(c *fiber.Ctx) error {
	id := c.Params("id")

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("DELETE FROM city WHERE id = ?", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "Kota tidak ditemukan"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
