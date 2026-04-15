package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	// Baca konfigurasi dari env variables
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")

	// Default values jika env var tidak ada
	if dbUser == "" {
		dbUser = "root"
	}
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbName == "" {
		dbName = "db_ticketing"
	}

	// Pastikan server MySQL terhubung dan database ada.
	adminDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true", dbUser, dbPassword, dbHost, dbPort)
	adminDB, err := sql.Open("mysql", adminDSN)
	if err != nil {
		log.Fatal("Gagal koneksi server MySQL: ", err)
	}
	defer adminDB.Close()

	if err := adminDB.Ping(); err != nil {
		log.Fatal("MySQL tidak bisa diakses: ", err)
	}

	safeDBName := strings.ReplaceAll(dbName, "`", "``")
	if _, err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", safeDBName)); err != nil {
		log.Fatal("Gagal membuat database: ", err)
	}

	// Buat DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Gagal koneksi ke database: ", err)
	}

	return db
}
