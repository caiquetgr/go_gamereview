package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	rtr := gin.Default()
	rtr.GET("/games", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"id": "1",
		})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: rtr,
	}

	srvErrors := make(chan error, 1)

	go func() {
		log.Println("server starting up")
		srvErrors <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL)

	select {
	case err := <-srvErrors:
		log.Fatal(fmt.Errorf("server error: %w", err))
		os.Exit(1)

	case sig := <-quit:
		log.Println("Server shutting down with signal", sig)

		ctx, canc := context.WithTimeout(context.Background(), 5*time.Second)
		defer canc()

		if err := srv.Shutdown(ctx); err != nil {
			srv.Close()
			log.Fatal("Server shutdown:", err)
		}

		select {
		case <-ctx.Done():
			log.Println("server shutdown timed out")
		}

		log.Println("server existing")
	}
}
