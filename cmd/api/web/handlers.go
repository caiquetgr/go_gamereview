package web

import (
	"net/http"

	v1 "github.com/caiquetgr/go_gamereview/cmd/api/web/v1"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/uptrace/bun"

	"github.com/gin-gonic/gin"
)

type ApiConfig struct {
	DB            *bun.DB
	KafkaProducer *kafka.Producer
}

func Handlers(cfg ApiConfig) http.Handler {
	h := gin.Default()
	v1.Handle(h, v1.ApiConfig{
		DB:            cfg.DB,
		KafkaProducer: cfg.KafkaProducer,
	})
	return h
}
