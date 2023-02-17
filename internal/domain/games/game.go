package games

import (
	"context"
)

type Repository interface {
	FindAll(ctx context.Context, page int, pageSize int) ([]Game, bool, error)
}

type GameService struct {
	repo Repository
}

func NewGameService(repo Repository) GameService {
	return GameService{
		repo: repo,
	}
}

func (s GameService) GetAllGames(ctx context.Context, page int, pageSize int) ([]Game, bool, error) {
	games, hasNext, err := s.repo.FindAll(ctx, page, pageSize)

	if err != nil {
		return nil, false, err
	}
	return games, hasNext, err
}
