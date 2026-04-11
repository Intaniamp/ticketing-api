package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	// Pastikan nama database di Laragon kamu adalah 'db_ticketing'
	dsn := "root:@tcp(localhost:3306)/db_ticketing?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi: ", err)
	}
	return db
}
