package games

import (
	"fmt"
	"testing"

	"github.com/caiquetgr/go_gamereview/internal/platform/database/dbtest"
)

func TestMain(m *testing.M) {
	c, err := dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)
	m.Run()
}
