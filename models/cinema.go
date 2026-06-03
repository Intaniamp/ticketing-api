package models

// Cinema merepresentasikan balikan data utuh dari database
type Cinema struct {
	ID         int    `json:"id" example:"1"`
	CinemaName string `json:"cinema_name" example:"Beachwalk XXI"`
	CityID     int    `json:"city_id" example:"1"`
	Address    string `json:"address" example:"Beachwalk Shopping Center Lt. 2, Jl. Pantai Kuta"`
}

// CinemaRequest digunakan khusus untuk menangkap input create/update dari user (tanpa ID)
type CinemaRequest struct {
	CinemaName string `json:"cinema_name" example:"Beachwalk XXI"`
	CityID     int    `json:"city_id" example:"1"`
	Address    string `json:"address" example:"Beachwalk Shopping Center Lt. 2, Jl. Pantai Kuta"`
}
