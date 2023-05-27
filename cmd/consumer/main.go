package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/caiquetgr/go_gamereview/cmd/consumer/kafka"
	"github.com/caiquetgr/go_gamereview/internal/platform/database"
	kafkafactory "github.com/caiquetgr/go_gamereview/internal/platform/events/kafka"
)

func main() {
	ctx := context.Background()

	db := database.OpenConnection(database.DbConfig{
		Host:            "localhost:5432",
		User:            "postgres",
		Password:        "postgres",
		Database:        "gamereview",
		ApplicationName: "go_gamereview",
	})
	defer db.Close()

	err := database.Migrate(ctx, db)
	if err != nil {
		panic(err)
	}

	kp := kafkafactory.CreateKafkaProducer(kafkafactory.ProducerConfig{
		BootstrapServers: "localhost:9092",
		Acks:             "all",
	})
	defer kp.Close()

	kcCreator := kafkafactory.CreateKafkaConsumer

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	kafka.Handle(kafka.KafkaHandlerConfig{
		DB:                  db,
		KafkaProducer:       kp,
		KafkaConsumerCreate: kcCreator,
		SigChan:             sigchan,
	})
}
