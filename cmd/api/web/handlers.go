package web

import (
	v1 "github.com/caiquetgr/go_gamereview/cmd/api/web/v1"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handlers() http.Handler {
	h := gin.Default()
	v1.Handle(h)
	return h
}
