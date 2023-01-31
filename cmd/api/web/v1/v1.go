package v1

import (
	"github.com/caiquetgr/go_gamereview/cmd/api/web/v1/handler"
	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/gin-gonic/gin"
)

func Handle(ge *gin.Engine) {
	g := ge.Group("/v1")

	gh := handler.NewGameHandler(games.NewGameService())

	g.GET("/games", gh.GetAll)
}
