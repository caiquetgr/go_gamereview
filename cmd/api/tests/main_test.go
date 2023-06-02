package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/caiquetgr/go_gamereview/foundation/test"
	"github.com/testcontainers/testcontainers-go/modules/compose"
)

const (
	HttpServerPort    = ":8080"
	HttpServerAddress = "http://localhost"
)

var comp compose.ComposeStack

func TestMain(m *testing.M) {
	ctx := context.Background()
	comp, err := test.InitDependencies(ctx)
	if err != nil {
		panic(err)
	}

	m.Run()

	err = comp.Down(ctx)
	if err != nil {
		panic(err)
	}
}

func getServerUrl() string {
	return fmt.Sprintf("%s%s", HttpServerAddress, HttpServerPort)
}

func TestMain_IsTestWorking(t *testing.T) {
	t.Log("starting test")
	time.Sleep(1 * time.Second)
	t.Log("finished test")
}
