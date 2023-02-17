package web

import (
	"net/http"

	v1 "github.com/caiquetgr/go_gamereview/cmd/api/web/v1"
	"github.com/uptrace/bun"

	"github.com/gin-gonic/gin"
)

type ApiConfig struct {
	DB *bun.DB
}

func Handlers(cfg ApiConfig) http.Handler {
	h := gin.Default()
	v1.Handle(h, v1.ApiConfig{
		DB: cfg.DB,
	})
	return h
}
