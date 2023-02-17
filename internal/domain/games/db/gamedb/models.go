package gamedb

import (
	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/caiquetgr/go_gamereview/internal/platform/database"
	"github.com/uptrace/bun"
)

type GameDbModel struct {
	bun.BaseModel `bun:"table:games"`
	Base          database.BaseDbModel `bun:",extend"`
	Name          string               `bun:"name,notnull"`
	Year          int                  `bun:"year,notnull"`
	Platform      string               `bun:"platform,notnull"`
	Genre         string               `bun:"genre,notnull"`
	Publisher     string               `bun:"publisher,notnull"`
}

func toGame(g GameDbModel) games.Game {
	return games.Game{
		ID:        g.Base.ID,
		Name:      g.Name,
		Year:      g.Year,
		Platform:  g.Platform,
		Genre:     g.Genre,
		Publisher: g.Publisher,
	}
}

func toGameDb(g games.Game) GameDbModel {
	return GameDbModel{
		Name:      g.Name,
		Year:      g.Year,
		Platform:  g.Platform,
		Genre:     g.Genre,
		Publisher: g.Publisher,
	}
}
