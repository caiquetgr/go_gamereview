package games

import (
	"context"
)

type GameService struct {
}

func NewGameService() *GameService {
	return &GameService{}
}

func (s *GameService) GetAllGames(ctx context.Context) []*Game {
	return []*Game{
		{
			Id:        1,
			Name:      "Teste",
			Year:      2022,
			Platform:  "Super Nintendo",
			Genre:     "Platform",
			Publisher: "Nintendo",
		},
	}
}
