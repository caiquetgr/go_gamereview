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

	"github.com/caiquetgr/go_gamereview/cmd/api/web"
	"github.com/caiquetgr/go_gamereview/internal/platform/database"
	"github.com/caiquetgr/go_gamereview/internal/platform/events/kafka"
)

type DbConfig struct {
	Host            string
	User            string
	Password        string
	Database        string
	ApplicationName string
}

type KafkaProducerConfig struct {
	BootstrapServers string
	Acks             string
}

type HttpServerConfig struct {
	Addr string
}

type AppConfig struct {
	DbConfig         DbConfig
	KPConfig         KafkaProducerConfig
	HttpServerConfig HttpServerConfig
	AppReadyChan     chan struct{}
}

func main() {
	appConfig := AppConfig{
		DbConfig: DbConfig{
			Host:            "localhost:5432",
			User:            "postgres",
			Password:        "postgres",
			Database:        "gamereview",
			ApplicationName: "go_gamereview",
		},
		KPConfig: KafkaProducerConfig{
			BootstrapServers: "localhost:9092",
			Acks:             "all",
		},
		HttpServerConfig: HttpServerConfig{
			Addr: ":8080",
		},
		AppReadyChan: make(chan struct{}),
	}

	Run(context.Background(), appConfig)
}

func Run(ctx context.Context, appConfig AppConfig) {
	db := database.OpenConnection(database.DbConfig{
		Host:            appConfig.DbConfig.Host,
		User:            appConfig.DbConfig.User,
		Password:        appConfig.DbConfig.Password,
		Database:        appConfig.DbConfig.Database,
		ApplicationName: appConfig.DbConfig.ApplicationName,
	})
	defer db.Close()

	err := database.Migrate(ctx, db)
	if err != nil {
		panic(err)
	}

	kp := kafka.CreateKafkaProducer(kafka.ProducerConfig{
		BootstrapServers: appConfig.KPConfig.BootstrapServers,
		Acks:             appConfig.KPConfig.Acks,
	})
	defer kp.Close()

	rtr := web.Handlers(web.ApiConfig{
		DB:            db,
		KafkaProducer: kp,
	})

	srv := &http.Server{
		Addr:    appConfig.HttpServerConfig.Addr,
		Handler: rtr,
	}

	srvErrors := make(chan error, 1)

	go func() {
		log.Println("server started")
		srvErrors <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	close(appConfig.AppReadyChan)

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

		<-ctx.Done()
		log.Println("server shutdown timed out")
		log.Println("server exiting")
	}
}
