package models

type Film struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Duration    int     `json:"duration"`
	Synopsis    string  `json:"synopsis"`
	PosterURL   string  `json:"poster_url"`
	AgeRating   string  `json:"age_rating"`
	ReleaseYear int     `json:"release_year"`
	Genres      []Genre `json:"genres"`
}

// untuk request body create/update film
type FilmRequest struct {
	Title       string `json:"title" example:"Scream 7"`
	Duration    int    `json:"duration" example:"114"`
	Synopsis    string `json:"synopsis" example:"Ghostface is back in town..."`
	AgeRating   string `json:"age_rating" example:"16+"`
	ReleaseYear int    `json:"release_year" example:"2026"`
	GenreIDs    []int  `json:"genre_ids" example:"1,3"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
