package test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/caiquetgr/go_gamereview/cmd/api/config"
	"github.com/caiquetgr/go_gamereview/internal/platform/database"
	k "github.com/caiquetgr/go_gamereview/internal/platform/events/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/uptrace/bun"
)

type AppIntegrationTest struct {
	Db        *bun.DB
	Kp        *kafka.Producer
	Kc        *kafka.Consumer
	KcCreator func(kcc k.ConsumerConfig) *kafka.Consumer
	Teardown  func()
}

func BuildAppConfig(ctx context.Context) config.AppConfig {
	return config.AppConfig{
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
}

func NewIntegrationTest(ctx context.Context) AppIntegrationTest {
	cfg := BuildAppConfig(ctx)

	db := database.OpenConnection(database.DbConfig{
		Host:            cfg.DbConfig.Host,
		User:            cfg.DbConfig.User,
		Password:        cfg.DbConfig.Password,
		Database:        cfg.DbConfig.Database,
		ApplicationName: cfg.DbConfig.ApplicationName,
	})

	err := database.Migrate(ctx, db)
	if err != nil {
		panic(err)
	}

	kp := k.CreateKafkaProducer(k.ProducerConfig{
		BootstrapServers: cfg.KPConfig.BootstrapServers,
		Acks:             cfg.KPConfig.Acks,
	})

	kc := k.CreateKafkaConsumer(k.ConsumerConfig{
		BootstrapServers: cfg.KPConfig.BootstrapServers,
		GroupId:          "go_gamereview",
		AutoOffsetReset:  "earliest",
	})

	return AppIntegrationTest{
		Db:        db,
		Kp:        kp,
		Kc:        kc,
		KcCreator: k.CreateKafkaConsumer,
		Teardown: func() {
			fmt.Println("Tearing down integration test")
			db.Close()
			kp.Close()
			kc.Close()
		},
	}
}

func WaitUntil(f func() bool, timeout time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	defer cancel()

	for {
		select {
		case <-ticker.C:
			log.Printf("trying to find")
			if f() {
				return true, nil
			}
		case <-ctx.Done():
			log.Printf("timeout")
			return false, fmt.Errorf("timed out waiting until func return true")
		}
	}
}
