package models

// Studio merepresentasikan balikan data utuh dari database
type Studio struct {
	ID        int    `json:"id" example:"32"`
	CinemaID  int    `json:"cinema_id" example:"7"`
	StudioName string `json:"studio_name" example:"Regular"`
}

// StudioRequest digunakan khusus untuk menangkap input create/update dari user (tanpa ID)
type StudioRequest struct {
	CinemaID  int    `json:"cinema_id" example:"7"`
	StudioName string `json:"studio_name" example:"Regular"`
}
