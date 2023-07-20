package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	handler "github.com/caiquetgr/go_gamereview/cmd/consumer/kafka/v1"
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

type KafkaMessageHandler interface {
	HandleMessage(c *kafka.Consumer, m *kafka.Message)
	GetTopic() string
}

type KafkaHandlerConfig struct {
	DB                  *bun.DB
	KafkaProducer       *kafka.Producer
	KafkaConsumerCreate func(kcc k.ConsumerConfig) *kafka.Consumer
	SigChan             <-chan (os.Signal)
	StopChan            <-chan (struct{})
}

type StartHandler struct {
	waitGroupDone   func()
	ctx             context.Context
	consumerCreator func() *kafka.Consumer
	handler         KafkaMessageHandler
}

func Handle(cfg KafkaHandlerConfig) {
	gs := games.NewGameService(
		gamedb.NewGameRepositoryBun(cfg.DB),
		event.NewGameEventProducer(NewGameTopic, cfg.KafkaProducer),
	)

	handlers := []KafkaMessageHandler{
		handler.BuildNewGameEventHandler(gs),
	}

	var wg sync.WaitGroup
	wg.Add(len(handlers))

	f := func() *kafka.Consumer {
		return cfg.KafkaConsumerCreate(k.ConsumerConfig{
			BootstrapServers: "localhost:9092",
			GroupId:          "go_gamereview",
			AutoOffsetReset:  "earliest",
		})
	}

	ctx, cancel := context.WithCancel(context.Background())

	for _, h := range handlers {
		sh := StartHandler{
			handler:         h,
			waitGroupDone:   func() { wg.Done() },
			consumerCreator: f,
			ctx:             ctx,
		}
		go startHandler(sh)
	}

	for run := true; run; {
		select {
		case sig := <-cfg.SigChan:
			log.Println("stopping kafka listeners with signal:", sig)
			run = false
		case <-cfg.StopChan:
			log.Println("stopping kafka listeners with stop channel read")
			run = false
		}
	}

	cancel()
	wg.Wait()
}

func startHandler(s StartHandler) {
	defer s.waitGroupDone()

	c := s.consumerCreator()
	defer c.Close()

	h := s.handler
	t := h.GetTopic()

	err := c.Subscribe(NewGameTopic, nil)
	if err != nil {
		panic(err)
	}

	log.Printf("starting consumer for topic %v\n", t)

	for run := true; run; {
		select {
		case <-s.ctx.Done():
			log.Printf("stopping consumer for topic %v by context cancellation\n", t)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}
			log.Println(ev)

			switch e := ev.(type) {
			case *kafka.Message:
				log.Printf("-- Message on %s: %s\n", e.TopicPartition, string(e.Value))
				log.Printf("-- Headers: %s\n", e.Headers)
				h.HandleMessage(c, e)

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
