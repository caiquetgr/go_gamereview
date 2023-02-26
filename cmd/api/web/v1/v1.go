package v1

import (
	"github.com/caiquetgr/go_gamereview/cmd/api/web/v1/handler"
	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/caiquetgr/go_gamereview/internal/domain/games/db/gamedb"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type ApiConfig struct {
	DB *bun.DB
}

func Handle(ge *gin.Engine, cfg ApiConfig) {
	g := ge.Group("/v1")

	gh := handler.NewGameHandler(
		games.NewGameService(gamedb.NewGameRepositoryBun(cfg.DB)),
	)

	g.GET("/games", gh.GetAll)
	g.POST("/games", gh.Create)
}
