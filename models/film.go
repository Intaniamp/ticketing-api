package models

type Film struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	Rating   string `json:"rating"`
	Synopsis string `json:"synopsis"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
