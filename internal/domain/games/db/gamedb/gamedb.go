package gamedb

import (
	"context"

	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/uptrace/bun"
)

type GameRepository interface {
	FindAll(ctx context.Context) []games.Game
}

type GameRepositoryBun struct {
	db *bun.DB
}

func NewGameRepositoryBun(DB *bun.DB) GameRepositoryBun {
	return GameRepositoryBun{
		db: DB,
	}
}

func (gr GameRepositoryBun) FindAll(ctx context.Context) []games.Game {
	var games []games.Game
	err := gr.db.NewSelect().Model(&games).Limit(20).Scan(ctx)
	if err != nil {
		panic(err)
	}
	return games
}
