package models

// Schedule merepresentasikan balikan data utuh dari database
type Schedule struct {
	ID       int    `json:"id" example:"1"`
	FilmID   int    `json:"film_id" example:"1"`
	StudioID int    `json:"studio_id" example:"32"`
	Date     string `json:"date" example:"2026-06-03"`
	Time     string `json:"time" example:"13:00:00"`
	Price    float64 `json:"price" example:"40000.00"`
}

// ScheduleRequest digunakan khusus untuk menangkap input create/update dari user (tanpa ID)
type ScheduleRequest struct {
	FilmID   int     `json:"film_id" example:"1"`
	StudioID int     `json:"studio_id" example:"32"`
	Date     string  `json:"date" example:"2026-06-03"`
	Time     string  `json:"time" example:"13:00:00"`
	Price    float64 `json:"price" example:"40000.00"`
}
