package models

type Seat struct {
	ID         int    `json:"id" example:"1"`
	StudioID   int    `json:"studio_id" example:"1"`
	SeatNumber string `json:"seat_number" example:"A1"`
}

// untuk input banyak kursi sekaligus
type BulkSeatRequest struct {
	StudioID    int      `json:"studio_id" example:"1"`
	SeatNumbers []string `json:"seat_numbers" example:"A1,A2,A3,B1,B2,B3"`
}
