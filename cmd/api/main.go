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

	"github.com/caiquetgr/go_gamereview/cmd/api/config"
	"github.com/caiquetgr/go_gamereview/cmd/api/web"
	"github.com/caiquetgr/go_gamereview/internal/platform/database"
	"github.com/caiquetgr/go_gamereview/internal/platform/events/kafka"
)

func main() {
	appConfig := config.AppConfig{
		DbConfig: config.DbConfig{
			Host:            "localhost:5432",
			User:            "postgres",
			Password:        "postgres",
			Database:        "gamereview",
			ApplicationName: "go_gamereview",
		},
		KPConfig: config.KafkaProducerConfig{
			BootstrapServers: "localhost:9092",
			Acks:             "all",
		},
		HttpServerConfig: config.HttpServerConfig{
			Addr: ":8080",
		},
	}

	Run(context.Background(), appConfig)
}

func Run(ctx context.Context, cfg config.AppConfig) {
	db := database.OpenConnection(database.DbConfig{
		Host:            cfg.DbConfig.Host,
		User:            cfg.DbConfig.User,
		Password:        cfg.DbConfig.Password,
		Database:        cfg.DbConfig.Database,
		ApplicationName: cfg.DbConfig.ApplicationName,
	})
	defer db.Close()

	err := database.Migrate(ctx, db)
	if err != nil {
		panic(err)
	}

	kp := kafka.CreateKafkaProducer(kafka.ProducerConfig{
		BootstrapServers: cfg.KPConfig.BootstrapServers,
		Acks:             cfg.KPConfig.Acks,
	})
	defer kp.Close()

	rtr := web.Handlers(web.ApiConfig{
		DB:            db,
		KafkaProducer: kp,
	})

	srv := &http.Server{
		Addr:    cfg.HttpServerConfig.Addr,
		Handler: rtr,
	}

	srvErrors := make(chan error, 1)

	go func() {
		log.Println("server started")
		srvErrors <- srv.ListenAndServe()
	}()

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-srvErrors:
		log.Fatal(fmt.Errorf("server error: %w", err))

	case sig := <-quitSignal:
		log.Println("Server shutting down with signal", sig)
		stop(srv)
	}
}

func stop(srv *http.Server) {
	ctx, canc := context.WithTimeout(context.Background(), 5*time.Second)
	defer canc()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
		srv.Close()
	}

	<-ctx.Done()
	log.Println("server finished")
}
