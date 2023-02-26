package handler

import (
	"net/http"

	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/caiquetgr/go_gamereview/internal/platform/web"
	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gs games.GameService
}

func NewGameHandler(gs games.GameService) GameHandler {
	return GameHandler{gs: gs}
}

func (h GameHandler) GetAll(c *gin.Context) {
	page, pageSize, err := web.QueryPageParams(c.Request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	games, hasNext, err := h.gs.GetAllGames(c.Request.Context(), page, pageSize)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"games":   games,
			"hasNext": hasNext,
		})
	}
}

func (h GameHandler) Create(c *gin.Context) {
	var newGame games.NewGame

	if err := c.ShouldBindJSON(&newGame); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad request": err.Error()})
		return
	}

	game, err := h.gs.CreateGame(c.Request.Context(), newGame)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusCreated, game)
	}
}
