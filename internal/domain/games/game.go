package games

import (
	"context"
)

type GameService struct {
	repo GameRepository
}

func NewGameService(repo GameRepository) GameService {
	return GameService{
		repo: repo,
	}
}

func (s GameService) GetAllGames(ctx context.Context) []Game {
	return s.repo.FindAll(ctx)
}
