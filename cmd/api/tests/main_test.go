package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/caiquetgr/go_gamereview/foundation/test"
	"github.com/testcontainers/testcontainers-go/modules/compose"
)

var comp compose.ComposeStack

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	comp, err = test.InitDependencies(ctx)

	if err != nil {
		panic(err)
	}

	defer func() {
		err := comp.Down(ctx)
		if err != nil {
			fmt.Printf("Failed to shutdown docker-compose: %v", err)
		}
	}()

	m.Run()
}
