package main

import (
	"context"
	"fmt"
	"github.com/caiquetgr/go_gamereview/cmd/api/web"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rtr := web.Handlers()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: rtr,
	}

	srvErrors := make(chan error, 1)

	go func() {
		log.Println("server started")
		srvErrors <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL)

	select {
	case err := <-srvErrors:
		log.Fatal(fmt.Errorf("server error: %w", err))

	case sig := <-quit:
		log.Println("Server shutting down with signal", sig)

		ctx, canc := context.WithTimeout(context.Background(), 5*time.Second)
		defer canc()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server shutdown:", err)
			srv.Close()
		}

		select {
		case <-ctx.Done():
			log.Println("server shutdown timed out")
		}

		log.Println("server existing")
	}
}
