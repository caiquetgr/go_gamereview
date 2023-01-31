package games

import (
	"github.com/caiquetgr/go_gamereview/internal/platform/database"
	"github.com/uptrace/bun"
)

type Game struct {
	ID        int64
	Name      string
	Year      int
	Platform  string
	Genre     string
	Publisher string
}

type GameDbModel struct {
	bun.BaseModel `bun:"table:games"`
	Base          database.BaseDbModel `bun:",extend"`
	Name          string               `bun:"name,notnull"`
	Year          int                  `bun:"year,notnull"`
	Platform      string               `bun:"platform,notnull"`
	Genre         string               `bun:"genre,notnull"`
	Publisher     string               `bun:"publisher,notnull"`
}

func (g *Game) NewGameDbModel() *GameDbModel {
	return &GameDbModel{
		Name:      g.Name,
		Year:      g.Year,
		Platform:  g.Platform,
		Genre:     g.Genre,
		Publisher: g.Publisher,
	}
}

func (g *GameDbModel) toGame() *Game {
	return &Game{
		ID:        g.Base.ID,
		Name:      g.Name,
		Year:      g.Year,
		Platform:  g.Platform,
		Genre:     g.Genre,
		Publisher: g.Publisher,
	}
}
