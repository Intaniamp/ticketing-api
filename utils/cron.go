package utils

import (
	"fmt"
	"ticketing-api/config"
	"time"
)

func StartTicketSweeper() {
	go func() {
		for {
			db := config.ConnectDB()

			rows, err := db.Query("SELECT user_id FROM booking WHERE status = 'pending' AND expired_at <= NOW()")
			var expiredUserIDs []string

			if err == nil {
				for rows.Next() {
					var uid string
					// Scan ID user dan masukkan ke dalam list
					if err := rows.Scan(&uid); err == nil {
						expiredUserIDs = append(expiredUserIDs, uid)
					}
				}
				rows.Close()
			}

			query := "UPDATE booking SET status = 'expired' WHERE status = 'pending' AND expired_at <= NOW()"
			result, err := db.Exec(query)

			if err == nil {
				affected, _ := result.RowsAffected()
				if affected > 0 {
					fmt.Printf("[SYSTEM] 🧹 Membersihkan %d tiket yang expired otomatis.\n", affected)

					for _, userID := range expiredUserIDs {
						SendPaymentNotification(userID, "Waktu pembayaran habis! Booking Anda telah dibatalkan.", "timeout")
					}
				}
			}

			db.Close()

			time.Sleep(1 * time.Minute)
		}
	}()
}
