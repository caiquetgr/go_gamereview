package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

const (
	dcFile = "docker-compose-test.yml"
)

func InitDependencies(t *testing.T) *tc.dockerCompose {
	c, err := tc.NewDockerCompose(dcFile)
	assert.NoError(t, err, "NewDockerComposeAPI()")
	return c
}
