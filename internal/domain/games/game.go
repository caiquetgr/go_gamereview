package games

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindAll(ctx context.Context, page int, pageSize int) ([]Game, bool, error)
	Create(ctx context.Context, g *Game) (*Game, error)
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
	// TODO: create page struct
	games, hasNext, err := s.repo.FindAll(ctx, page, pageSize)

	if err != nil {
		return nil, false, err
	}

	// TODO: map games return with camel or snake case
	return games, hasNext, err
}

func (s GameService) CreateGame(ctx context.Context, ng NewGame) (*Game, error) {
	g := Game{
		ID:        uuid.New(),
		Name:      ng.Name,
		Year:      ng.Year,
		Platform:  ng.Platform,
		Genre:     ng.Genre,
		Publisher: ng.Publisher,
	}

	return s.repo.Create(ctx, &g)
}
