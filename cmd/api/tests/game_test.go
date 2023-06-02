package tests

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/caiquetgr/go_gamereview/cmd/api/config"
	"github.com/caiquetgr/go_gamereview/cmd/api/web"
	"github.com/caiquetgr/go_gamereview/foundation/test"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/uptrace/bun"
)

type GameTest struct {
	app http.Handler
}

func TestGames(t *testing.T) {
	cfg := BuildAppConfig(context.Background(), comp)
	web.Handlers(web.ApiConfig{
		DB:            &bun.DB{},
		KafkaProducer: &kafka.Producer{},
	})
}

func BuildAppConfig(ctx context.Context, comp compose.ComposeStack) config.AppConfig {
	containers, err := test.GetContainers(ctx, comp)
	if err != nil {
		panic(err)
	}

	dbContainer := containers[test.DatabaseService]
	kafkaContainer := containers[test.KafkaService]

	dbAddress := test.GetContainerAddress(ctx, dbContainer, "5432")
	kafkaAddress := test.GetContainerAddress(ctx, kafkaContainer, "9092")

	return config.AppConfig{
		DbConfig: config.DbConfig{
			Host:            dbAddress,
			User:            "postgres",
			Password:        "postgres",
			Database:        "gamereview",
			ApplicationName: "go_gamereview",
		},
		KPConfig: config.KafkaProducerConfig{
			BootstrapServers: kafkaAddress,
			Acks:             "all",
		},
		HttpServerConfig: config.HttpServerConfig{
			Addr: HttpServerPort,
		},
	}
}

func TestGetGames(t *testing.T) {
	res, err := http.Get(fmt.Sprintf("%s/v1/games", getServerUrl()))
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	games, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Logf("%s", games)
}
