package games

import (
	"testing"
	"time"

	"github.com/caiquetgr/go_gamereview/foundation/test"
)

func Test_Game(t *testing.T) {
	t.Log("running!")
	test.InitDependencies(t)
	t.Log("test running!")
	time.Sleep(10 * time.Second)
	t.Log("test ended!")
}
