package main

import (
	"context"
	"testing"
	"time"

	"github.com/caiquetgr/go_gamereview/foundation/test"
	"github.com/testcontainers/testcontainers-go/modules/compose"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	comp, err := test.InitDependencies(ctx)
	defer comp.Down(ctx)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	appConfig := buildAppConfig(ctx, comp)
	ctxTimeout, canc := context.WithTimeout(context.Background(), 5*time.Second)
	defer canc()

	go func() {
		Run(ctx, appConfig)
	}()
}

func buildAppConfig(ctx context.Context, comp compose.ComposeStack) AppConfig {
	containers, err := test.GetContainers(ctx, comp)
	if err != nil {
		panic(err)
	}

	dbContainer := containers[test.DatabaseService]
	kafkaContainer := containers[test.KafkaService]

	dbAddress := test.GetContainerAddress(ctx, dbContainer, "5432")
	kafkaAddress := test.GetContainerAddress(ctx, kafkaContainer, "9092")

	return AppConfig{
		DbConfig: DbConfig{
			Host:            dbAddress,
			User:            "postgres",
			Password:        "postgres",
			Database:        "gamereview",
			ApplicationName: "go_gamereview",
		},
		KPConfig: KafkaProducerConfig{
			BootstrapServers: kafkaAddress,
			Acks:             "all",
		},
		HttpServerConfig: HttpServerConfig{
			Addr: ":8080",
		},
		AppReadyChan: make(chan struct{}),
	}
}
