package models

// City merepresentasikan balikan data utuh dari database
type City struct {
	ID       int    `json:"id" example:"1"`
	CityName string `json:"city_name" example:"Jakarta"`
}

// CityRequest digunakan khusus untuk menangkap input create/update dari user (tanpa ID)
type CityRequest struct {
	CityName string `json:"city_name" example:"Jakarta"`
}
