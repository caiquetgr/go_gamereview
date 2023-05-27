package test

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

const (
	dcFile          = "../../../docker-compose-test.yml"
	DatabaseService = "db"
	KafkaService    = "kafka"
)

func InitDependencies(ctx context.Context) (tc.ComposeStack, error) {
	comp, err := tc.NewDockerCompose(dcFile)
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
