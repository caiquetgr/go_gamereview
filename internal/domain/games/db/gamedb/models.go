package gamedb

import (
	"time"

	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GameDbModel struct {
	bun.BaseModel `bun:"table:games"`
	ID            uuid.UUID `bun:"id,pk"`
	CreatedAt     time.Time `bun:",notnull,default:current_timestamp`
	ModifiedAt    time.Time `bun:",notnull,default:current_timestamp`
	Name          string    `bun:"name,notnull"`
	Year          int       `bun:"year,notnull"`
	Platform      string    `bun:"platform,notnull"`
	Genre         string    `bun:"genre,notnull"`
	Publisher     string    `bun:"publisher,notnull"`
}

func toGames(gs []GameDbModel) []games.Game {
	games := make([]games.Game, len(gs))
	for i, g := range gs {
		games[i] = toGame(g)
	}
	return games
}

func toGame(g GameDbModel) games.Game {
	return games.Game{
		ID:        g.ID,
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
