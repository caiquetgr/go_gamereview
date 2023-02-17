package games

import (
	"context"
)

type Repository interface {
	FindAll(ctx context.Context) []Game
}

type GameService struct {
	repo Repository
}

func NewGameService(repo Repository) GameService {
	return GameService{
		repo: repo,
	}
}

func (s GameService) GetAllGames(ctx context.Context) []Game {
	return s.repo.FindAll(ctx)
}
