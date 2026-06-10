package utils

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// 1. Buat struktur kartu pengenal untuk Client
type Client struct {
	Conn   *websocket.Conn
	UserID string // Kita simpan userID di sini agar tahu ini koneksi milik siapa
}

// 2. Ubah isi WebSocketManager agar menggunakan struct Client baru kita
type WebSocketManager struct {
	Clients   map[*Client]bool 
	Broadcast chan map[string]interface{}
	Mutex     sync.Mutex
}

// Inisialisasi global manager
var WSManager = WebSocketManager{
	Clients:   make(map[*Client]bool),
	Broadcast: make(chan map[string]interface{}),
}

// Jalankan ini di background (goroutine) untuk mendengarkan data broadcast
func StartWebSocketHub() {
	for {
		message := <-WSManager.Broadcast
		
		targetUserID, _ := message["user_id"].(string)

		clientsToNotify := make([]*Client, 0)

		WSManager.Mutex.Lock()
		for client := range WSManager.Clients {
			if client.UserID == targetUserID {
				clientsToNotify = append(clientsToNotify, client)
			}
		}
		WSManager.Mutex.Unlock() // Buka gembok dengan cepat!

		// Kirim pesan khusus ke client yang lolos filter saja
		for _, client := range clientsToNotify {
			err := client.Conn.WriteJSON(message) 
			if err != nil {
				log.Println("❌ Gagal kirim pesan ke client, menutup koneksi:", err)
				client.Conn.Close()
				
				WSManager.Mutex.Lock()
				delete(WSManager.Clients, client)
				WSManager.Mutex.Unlock()
			}
		}
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