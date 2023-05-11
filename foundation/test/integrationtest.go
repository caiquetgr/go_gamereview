package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

const (
	dcFile = "../../../docker-compose-test.yml"
)

func InitDependencies(ctx context.Context, t *testing.T) map[string]*testcontainers.DockerContainer {
	comp, err := tc.NewDockerCompose(dcFile)
	assert.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		assert.NoError(t, comp.Down(ctx), "compose.Down()")
	})

	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)

	assert.NoError(t, comp.Up(ctx, tc.Wait(true)), "compose.Up()")

	services := comp.Services()
	containers := make(map[string]*testcontainers.DockerContainer)

	for _, s := range services {
		c, err := comp.ServiceContainer(ctx, s)
		assert.NoError(t, err, "ServiceContainer()")
		containers[s] = c
	}

	return containers
}
