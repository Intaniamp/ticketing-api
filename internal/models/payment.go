package models

// Payment merepresentasikan tabel payment di database
type Payment struct {
	ID            string  `json:"id" example:"pay-uuid-123"`
	BookingID     string  `json:"booking_id" example:"booking-uuid-456"`
	Amount        float64 `json:"amount" example:"60000"`
	PaymentMethod string  `json:"payment_method" example:"OVO"`
	PaymentDate   string  `json:"payment_date" example:"2026-06-05 14:30:00"`
	PaymentStatus string  `json:"payment_status" example:"success"`
}

// PaymentRequest menangkap data saat user menekan tombol "Pay Now"
type PaymentRequest struct {
	BookingID     string `json:"booking_id" example:"booking-uuid-456"`
	PaymentMethod string `json:"payment_method" example:"Gopay"`
}
