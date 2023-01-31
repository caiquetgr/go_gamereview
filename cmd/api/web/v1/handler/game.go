package handler

import (
	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gs *games.GameService
}

func NewGameHandler(s *games.GameService) *GameHandler {
	return &GameHandler{
		gs: s,
	}
}

func (h *GameHandler) GetAll(c *gin.Context) {
	c.JSON(200, h.gs.GetAllGames(c.Request.Context()))
}
