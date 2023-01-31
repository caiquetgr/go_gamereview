package games

import "github.com/uptrace/bun"

type GameRepository interface {
	FindAll() []*Game
}

type GameRepositoryBun struct {
	db *bun.DB
}

func NewGameRepositoryBun(db *bun.DB) *GameRepositoryBun {
	return &GameRepositoryBun{
		db: db,
	}
}
