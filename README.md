# Cinema +

## 1. Project Overview
Cinema+ merupakan sebuah sistem multiplatform (APK dan Web) yang dirancang untuk memfasilitasi pemesanan tiket bioskop secara online. Sistem ini memungkinkan pengguna untuk melihat informasi film, memilih jadwal penayangan, menentukan kursi, serta melakukan proses pemesanan dan pembayaran tiket secara online.

Di dalam API ini telah mendukung pengelolaan data seperti film, genre, studio, dan jadwal, serta dilengkapi dengan mekanisme autentikasi berbasis JWT dan struktur basis data yang terstruktur untuk memastikan konsistensi serta kemudahan dalam pengembangan sistem.

## 2. Tech Stack
- library: Go (Golang) 1.26+
- Fiber v2 sebagai web framework
- MySQL sebagai database
- golang-jwt/jwt v5 untuk autentikasi JWT
- swaggo/swag dan fiber-swagger untuk dokumentasi API (Swagger UI)
- joho/godotenv untuk membaca konfigurasi environment dari file .env

## 3. Database Diagram
ERD sistem dapat dilihat pada gambar berikut:
[docs/erd.png](docs/erd.png)

## 4. Installation Guide
1. Clone repository.
2. Masuk ke folder project.
3. Install dependency Go:

	go mod tidy

4. Siapkan file environment .env di root project (contoh):

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

8. Server akan berjalan di:

	http://localhost:3000

## 5. API Documentation Link
- Swagger UI: http://localhost:3000/swagger/index.html
- Swagger JSON file: docs/swagger.json
- Swagger YAML file: docs/swagger.yaml