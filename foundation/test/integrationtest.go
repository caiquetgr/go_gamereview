package test

import (
	"context"
	"fmt"
	"time"

	"github.com/caiquetgr/go_gamereview/cmd/api/config"
	"github.com/caiquetgr/go_gamereview/internal/platform/database"
	k "github.com/caiquetgr/go_gamereview/internal/platform/events/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/uptrace/bun"
)

const (
	dcFile          = "../../../docker-compose-test.yml"
	DatabaseService = "db"
	KafkaService    = "kafka"
)

type AppIntegrationTest struct {
	Db        *bun.DB
	Kp        *kafka.Producer
	Kc        *kafka.Consumer
	KcCreator func(kcc k.ConsumerConfig) *kafka.Consumer
	Teardown  func()
}

func InitDependencies(ctx context.Context, composeFile string) (tc.ComposeStack, error) {
	comp, err := tc.NewDockerCompose(composeFile)
	if err != nil {
		return nil, err
	}

	err = comp.Up(ctx, tc.Wait(true))
	if err != nil {
		return nil, err
	}

	return comp, nil
}

func GetContainers(ctx context.Context, comp tc.ComposeStack) (map[string]*testcontainers.DockerContainer, error) {
	services := comp.Services()
	containers := make(map[string]*testcontainers.DockerContainer)

	for _, s := range services {
		c, err := comp.ServiceContainer(ctx, s)
		if err != nil {
			return nil, err
		}
		containers[s] = c
	}

	return containers, nil
}

func GetContainerAddress(ctx context.Context, c testcontainers.Container, containerPort string) string {
	host, _ := c.Host(ctx)
	port, _ := c.MappedPort(ctx, nat.Port(containerPort))
	return fmt.Sprintf("%s:%s", host, port.Port())
}

func BuildAppConfig(ctx context.Context, comp compose.ComposeStack) config.AppConfig {
	containers, err := GetContainers(ctx, comp)
	if err != nil {
		panic(err)
	}

	dbContainer := containers[DatabaseService]
	kafkaContainer := containers[KafkaService]

	dbAddress := GetContainerAddress(ctx, dbContainer, "5432")
	kafkaAddress := GetContainerAddress(ctx, kafkaContainer, "9092")

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
			Addr: ":8080",
		},
	}
}

func NewIntegrationTest(ctx context.Context, comp compose.ComposeStack) AppIntegrationTest {
	cfg := BuildAppConfig(ctx, comp)

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
		BootstrapServers: "localhost:9092",
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
			if f() {
				return true, nil
			}
		case <-ctx.Done():
			return false, fmt.Errorf("timed out waiting until func return true")
		}
	}
}
