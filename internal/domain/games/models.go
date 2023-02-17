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
