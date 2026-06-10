# Cinema +

## 1. Project Overview

Cinema+ merupakan sebuah sistem multiplatform (APK dan Web) yang dirancang untuk memfasilitasi pemesanan tiket bioskop secara online. Sistem ini memungkinkan pengguna untuk melihat informasi film, memilih jadwal penayangan, menentukan kursi, serta melakukan proses pemesanan dan pembayaran tiket secara online.

Di dalam API ini telah mendukung pengelolaan data seperti film, genre, studio, dan jadwal, serta dilengkapi dengan mekanisme autentikasi berbasis JWT dan struktur basis data yang terstruktur. Sistem ini juga dilengkapi dengan fitur real-time menggunakan WebSocket untuk notifikasi pembayaran dan timeout pesanan, serta Ticket Sweeper otomatis untuk membersihkan tiket yang kedaluwarsa guna memastikan konsistensi serta kemudahan dalam pengembangan sistem.

## 2. Tech Stack

* Bahasa: Go (Golang) 1.26+
* Framework: Fiber v2 sebagai web framework
* Database: MySQL
* Autentikasi: golang-jwt/jwt v5 untuk autentikasi JWT
* Real-time: [github.com/gofiber/websocket/v2] untuk fitur WebSocket
* Dokumentasi API: swaggo/swag dan fiber-swagger (Swagger UI)
* Environment: joho/godotenv untuk membaca konfigurasi dari file .env

## 3. Database Diagram

ERD sistem dapat dilihat pada gambar berikut:
[docs/erd.png]

## 4. Installation Guide

1. Clone repository.
2. Masuk ke folder project.
3. Install dependency Go:

go mod tidy

4. Siapkan file environment `.env` di root project (contoh):
JWT_SECRET=isi_secret
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=nama_database
DB_USERNAME=root
DB_PASSWORD=

5. Pastikan database MySQL sudah berjalan dan database sudah dibuat.
6. (Opsional) Generate ulang dokumen swagger:

swag init

7. Jalankan aplikasi:

go run main.go

8. **Testing API & Server:**
* Untuk testing di PC lokal (Web/Emulator), server berjalan di: `http://localhost:3000`
* Untuk testing menggunakan device HP, pastikan HP dan PC berada di jaringan WiFi yang sama. Gunakan IP lokal PC (contoh: `http://192.168.x.x:3000`) dan pastikan Windows Firewall diatur ke turn Off (atau izinkan untuk port 3000).

## 5. API Documentation Link

* Swagger UI: `http://localhost:3000/swagger/index.html`
* Swagger JSON file: `docs/swagger.json`
* Swagger YAML file: `docs/swagger.yaml`