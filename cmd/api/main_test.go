package main

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/caiquetgr/go_gamereview/foundation/test"
	"github.com/testcontainers/testcontainers-go/modules/compose"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	comp, err := test.InitDependencies(ctx)
	if err != nil {
		panic(err)
	}

	appConfig := buildAppConfig(ctx, comp)

	go func() {
		Run(ctx, appConfig)
	}()

	select {
	case <-time.After(10 * time.Second):
		panic("Timedout waiting App start")
	case <-appConfig.AppReadyChan:
		log.Println("App ready for integration tests")
	}

	m.Run()

	appConfig.AppStopChan <- struct{}{}

	select {
	case <-time.After(10 * time.Second):
		panic("Timedout waiting App stop")
	case <-appConfig.AppStopChan:
		log.Println("App shutdown for integration tests")
	}

	close(appConfig.AppStopChan)

	err = comp.Down(ctx)
	if err != nil {
		panic(err)
	}
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
		AppStopChan:  make(chan struct{}),
	}
}

func TestMain_IsTestWorking(t *testing.T) {
	t.Log("starting test")
	time.Sleep(10 * time.Second)
	t.Log("finished test")
}
