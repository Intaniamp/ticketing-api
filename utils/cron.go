package utils

import (
	"fmt"
	"ticketing-api/config"
	"time"
)

// StartTicketSweeper akan berjalan di background untuk membersihkan tiket yang expired
func StartTicketSweeper() {
	go func() {
		for {
			db := config.ConnectDB()

			// Ubah tiket 'pending' menjadi 'expired' jika waktu expired_at sudah terlewati
			query := "UPDATE booking SET status = 'expired' WHERE status = 'pending' AND expired_at <= NOW()"

			result, err := db.Exec(query)
			if err == nil {
				affected, _ := result.RowsAffected()
				if affected > 0 {
					fmt.Printf("[SYSTEM] 🧹 Membersihkan %d tiket yang expired otomatis.\n", affected)
				}
			}

			db.Close()

			// Suruh Golang tidur selama 1 menit sebelum mengecek ke database lagi
			time.Sleep(1 * time.Minute)
		}
	}()
}
