package utils

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Struktur untuk menyimpan client yang terkoneksi
type WebSocketManager struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan map[string]interface{}
	Mutex     sync.Mutex
}

// Inisialisasi global manager
var WSManager = WebSocketManager{
	Clients:   make(map[*websocket.Conn]bool),
	Broadcast: make(chan map[string]interface{}),
}

// Jalankan ini di background (goroutine) untuk dengerin data broadcast
func StartWebSocketHub() {
	for {
		message := <-WSManager.Broadcast
		WSManager.Mutex.Lock()
		
		// Kirim pesan ke SEMUA client yang lagi terkoneksi
		for client := range WSManager.Clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Println("❌ Gagal kirim pesan ke client, menutup koneksi:", err)
				client.Close()
				delete(WSManager.Clients, client)
			}
		}
		WSManager.Mutex.Unlock()
	}
}

// Fungsi pembantu untuk memicu notifikasi dari controller manapun
func SendPaymentNotification(userID string, message string, status string) {
	payload := map[string]interface{}{
		"type":    "PAYMENT_NOTIFICATION",
		"user_id": userID,
		"message": message,
		"status":  status,
	}
	// Masukkan ke channel broadcast
	WSManager.Broadcast <- payload
}