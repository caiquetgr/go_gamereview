package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/caiquetgr/go_gamereview/internal/domain/games/db/gamedb"
	"github.com/caiquetgr/go_gamereview/internal/domain/games/event"
	k "github.com/caiquetgr/go_gamereview/internal/platform/events/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/uptrace/bun"
)

const (
	NewGameTopic = "new-game-event"
)

type KafkaHandlerConfig struct {
	DB                  *bun.DB
	KafkaProducer       *kafka.Producer
	KafkaConsumerCreate func(kcc k.ConsumerConfig) *kafka.Consumer
	SigChan             <-chan (os.Signal)
}

func Handle(cfg KafkaHandlerConfig) {
	c := cfg.KafkaConsumerCreate(k.ConsumerConfig{
		BootstrapServers: "localhost:9092",
		GroupId:          "go_gamereview",
		AutoOffsetReset:  "earliest",
	})

	defer c.Close()

	err := c.SubscribeTopics([]string{NewGameTopic}, nil)
	if err != nil {
		panic(err)
	}

	gs := games.NewGameService(
		gamedb.NewGameRepositoryBun(cfg.DB),
		event.NewGameEventProducer(NewGameTopic, cfg.KafkaProducer),
	)

	run := true

	log.Println("starting consumer...")

	for run {
		select {
		case sig := <-cfg.SigChan:
			log.Println("stopping kafka listener with signal:", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				log.Printf("-- Message on %s: %s\n", e.TopicPartition, string(e.Value))
				log.Printf("-- Headers: %s\n", e.Headers)

				ng := &games.NewGame{}

				if err := json.Unmarshal(e.Value, ng); err != nil {
					fmt.Fprintf(os.Stderr, "Error unmarshalling message %v - error %v", e.Value, err)
					continue
				}

				game, err := gs.CreateGame(context.Background(), *ng)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error creating game: %v", err)
				} else {
					log.Println("created game", game)
				}
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "kafka error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				log.Println("ignored", e)
			}
		}
	}
}
