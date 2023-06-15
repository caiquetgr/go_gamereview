package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caiquetgr/go_gamereview/cmd/api/web"
	"github.com/caiquetgr/go_gamereview/foundation/test"
)

type GameTest struct {
	app http.Handler
}

func TestGames(t *testing.T) {
	ctx := context.Background()
	cfg := test.BuildAppConfig(ctx, comp)
	it := test.NewIntegrationTest(ctx, cfg)

	t.Cleanup(it.Teardown)

	tests := GameTest{
		app: web.Handlers(web.ApiConfig{
			DB:            it.Db,
			KafkaProducer: it.Kp,
		}),
	}

	t.Run("GetGamesList", tests.GetGames)
}

func (g GameTest) GetGames(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/v1/games", nil)
	w := httptest.NewRecorder()

	g.app.ServeHTTP(w, r)

	resp := map[string]interface{}{}

	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("[ERROR] Failed to decode response body: %v", err)
	}

	t.Logf("[INFO] Games: %v", resp)
}
