package games

import "github.com/google/uuid"

type Game struct {
	ID        uuid.UUID
	Name      string
	Year      int
	Platform  string
	Genre     string
	Publisher string
}

type NewGame struct {
	Name      string `json:"name"`
	Year      int    `json:"year"`
	Platform  string `json:"platform"`
	Genre     string `json:"genre"`
	Publisher string `json:"publisher"`
}
