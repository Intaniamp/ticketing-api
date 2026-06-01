package models

// Genre merepresentasikan balikan data utuh dari database
type Genre struct {
	ID        int    `json:"id" example:"1"`
	GenreName string `json:"genre_name" example:"Action"`
}

// GenreRequest digunakan khusus untuk menangkap input create/update dari user (tanpa ID)
type GenreRequest struct {
	GenreName string `json:"genre_name" example:"Action"`
}
