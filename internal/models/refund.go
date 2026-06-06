package models

// Refund merepresentasikan tabel refund di database
type Refund struct {
	ID           string  `json:"id" example:"ref-uuid-123"`
	BookingID    string  `json:"booking_id" example:"booking-uuid-456"`
	RefundDate   string  `json:"refund_date" example:"2026-06-05 15:00:00"`
	RefundAmount float64 `json:"refund_amount" example:"60000"`
}

// RefundRequest menangkap data saat user mengajukan refund
type RefundRequest struct {
	BookingID string `json:"booking_id" example:"booking-uuid-456"`
	Reason    string `json:"reason" example:"Salah pilih jam tayang"`
}
