package gamedb

import (
	"context"
	"fmt"

	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/uptrace/bun"
)

type GameRepositoryBun struct {
	db *bun.DB
}

func NewGameRepositoryBun(DB *bun.DB) GameRepositoryBun {
	return GameRepositoryBun{
		db: DB,
	}
}

func (gr GameRepositoryBun) FindAll(ctx context.Context, page int, pageSize int) ([]games.Game, bool, error) {
	var games []GameDbModel
	err := gr.db.NewSelect().
		Model(&games).
		Limit(pageSize + 1).
		Offset(pageSize * (page - 1)).
		Scan(ctx)

	if err != nil {
		return nil, false, fmt.Errorf("error finding games: %w", err)
	}

	hasNext := len(games) > pageSize

	if hasNext && len(games) > 0 {
		games = games[:len(games)-1]
	}

	return toGames(games), hasNext, nil
}

func (gr GameRepositoryBun) Create(ctx context.Context, g *games.Game) (*games.Game, error) {
	_, err := gr.db.NewInsert().Model(g).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return g, nil
}
