package web

import (
	v1 "github.com/caiquetgr/go_gamereview/cmd/api/web/v1"
	"github.com/uptrace/bun"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiConfig struct {
	DB *bun.DB
}

func Handlers(cfg ApiConfig) http.Handler {
	h := gin.Default()
	v1.Handle(h, cfg)
	return h
}
