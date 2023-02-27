package v1

import (
	"github.com/caiquetgr/go_gamereview/cmd/api/web/v1/handler"
	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/caiquetgr/go_gamereview/internal/domain/games/db/gamedb"
	"github.com/caiquetgr/go_gamereview/internal/domain/games/event"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type ApiConfig struct {
	DB            *bun.DB
	KafkaProducer *kafka.Producer
}

func Handle(ge *gin.Engine, cfg ApiConfig) {
	g := ge.Group("/v1")

	gh := handler.NewGameHandler(
		games.NewGameService(
			gamedb.NewGameRepositoryBun(cfg.DB),
			event.NewGameEventProducer("new-game-event", cfg.KafkaProducer),
		),
	)

	g.GET("/games", gh.GetAll)
	g.POST("/games", gh.CreateAsync)
}
