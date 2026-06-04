package models

// Booking merepresentasikan tabel booking untuk riwayat ringkas (History)
type Booking struct {
	ID          string  `json:"id" example:"98231a54-f58c-48a0-9c2b..."`
	UserID      string  `json:"user_id" example:"user-uuid-456"`
	ScheduleID  int     `json:"schedule_id" example:"1"`
	BookingDate string  `json:"booking_date" example:"2026-06-05 14:30:00"`
	ExpiredAt   string  `json:"expired_at" example:"2026-06-05 14:35:00"`
	Status      string  `json:"status" example:"pending"`
	TotalPrice  float64 `json:"total_price" example:"60000"`
}

// BookingRequest menangkap data saat user mau pesan tiket
type BookingRequest struct {
	ScheduleID int   `json:"schedule_id" example:"1"`
	SeatIDs    []int `json:"seat_ids" example:"15,16,17"`
}

// untuk halaman E-Ticket & Payment Details
type BookingDetailResponse struct {
	ID         string   `json:"id" example:"98231a54-f58c-48a0-9c2b..."`
	OrderID    string   `json:"order_id" example:"#CM-98231A54"`
	FilmTitle  string   `json:"film_title" example:"Scream 7"`
	PosterURL  string   `json:"poster_url" example:"uploads/posters/poster.jpg"`
	CinemaName string   `json:"cinema_name" example:"XX1 Living World"`
	Date       string   `json:"date" example:"2026-04-19"`
	Time       string   `json:"time" example:"17:30"`
	ExpiredAt  string   `json:"expired_at" example:"2026-06-05 14:35:00"`
	Status     string   `json:"status" example:"success"`
	TotalPrice float64  `json:"total_price" example:"60000"`
	Seats      []string `json:"seats" example:"F7,F8"`
}
