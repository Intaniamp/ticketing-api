package handlers

import (
	"database/sql"
	"fmt"
	"strings"
	"ticketing-api/config"
	"ticketing-api/models"

	"github.com/gofiber/fiber/v2"
)

type userResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type updateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Role  *string `json:"role"`
}

// GetAllUsers godoc
//
//	@Summary		Ambil semua user
//	@Description	Mengembalikan seluruh daftar user tanpa field password
//	@Tags			User
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		userResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/user [get]
func GetAllUsers(c *fiber.Ctx) error {
	db := config.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, email, role FROM user")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}
	defer rows.Close()

	var users []userResponse
	for rows.Next() {
		var u userResponse
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

// GetUserByID godoc
//
//	@Summary		Ambil user berdasarkan ID
//	@Description	Mengembalikan detail user tanpa field password berdasarkan ID
//	@Tags			User
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	userResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/user/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	db := config.ConnectDB()
	defer db.Close()

	var u userResponse
	err := db.QueryRow("SELECT id, name, email, role FROM user WHERE id = ?", id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "User tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(u)
}

// UpdateUser godoc
//
//	@Summary		Update user
//	@Description	Mengubah data user berdasarkan ID
//	@Tags			User
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int			true	"User ID"
//	@Param			request	body		models.User	true	"Data user"
//	@Success		200		{object}	userResponse
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/user/{id} [patch]
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var payload updateUserRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Invalid input"})
	}

	setClauses := make([]string, 0, 3)
	args := make([]interface{}, 0, 4)

	if payload.Name != nil {
		setClauses = append(setClauses, "name = ?")
		args = append(args, strings.TrimSpace(*payload.Name))
	}

	if payload.Email != nil {
		setClauses = append(setClauses, "email = ?")
		args = append(args, strings.TrimSpace(*payload.Email))
	}

	if payload.Role != nil {
		role := strings.TrimSpace(*payload.Role)
		if role != "user" && role != "admin" {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Role harus 'user' atau 'admin'"})
		}
		setClauses = append(setClauses, "role = ?")
		args = append(args, role)
	}

	if len(setClauses) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Message: "Tidak ada field yang diupdate"})
	}

	db := config.ConnectDB()
	defer db.Close()

	query := fmt.Sprintf("UPDATE user SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	args = append(args, id)

	result, err := db.Exec(
		query,
		args...,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "User tidak ditemukan"})
	}

	var updated userResponse
	err = db.QueryRow("SELECT id, name, email, role FROM user WHERE id = ?", id).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Email,
		&updated.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "User tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

// DeleteUser godoc
//
//	@Summary		Hapus user
//	@Description	Menghapus user berdasarkan ID
//	@Tags			User
//	@Security		BearerAuth
//	@Param			id	path		int	true	"User ID"
//	@Success		204	{string}	string
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/user/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	db := config.ConnectDB()
	defer db.Close()

	result, err := db.Exec("DELETE FROM user WHERE id = ?", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Message: err.Error()})
	}
	if affected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Message: "User tidak ditemukan"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
