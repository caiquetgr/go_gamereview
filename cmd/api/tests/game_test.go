package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caiquetgr/go_gamereview/cmd/api/web"
	"github.com/caiquetgr/go_gamereview/foundation/test"
	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/caiquetgr/go_gamereview/internal/domain/games/db/gamedb"
	"github.com/caiquetgr/go_gamereview/internal/domain/games/event"
	"github.com/stretchr/testify/assert"
)

type GameTest struct {
	app http.Handler
	gs  games.GameService
}

type GetGamesV1Response struct {
	Games   []games.Game
	HasNext bool
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
		gs: games.NewGameService(
			gamedb.NewGameRepositoryBun(it.Db),
			event.NewGameEventProducer("new-game-event", it.Kp),
		),
	}

	t.Run("GetGamesList", tests.GetGames)
}

func (g GameTest) GetGames(t *testing.T) {
	game := games.NewGame{
		Name:      "Super Ghouls'n Ghosts",
		Year:      1991,
		Platform:  "Super Nintendo",
		Genre:     "Platform",
		Publisher: "Capcom",
	}

	t.Log("\t Given a created game in database")
	{
		_, err := g.gs.CreateGame(context.Background(), game)
		if err != nil {
			t.Fatalf("\t Failed creating game in database for test: %v", err)
		}
	}

	resp := GetGamesV1Response{}

	t.Log("\t Should return 200 in a GET /v1/games request")
	{
		r := httptest.NewRequest(http.MethodGet, "/v1/games", nil)
		w := httptest.NewRecorder()

		g.app.ServeHTTP(w, r)

		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("[ERROR] Failed to decode response body: %v", err)
		}

		assert.Equal(t, http.StatusOK, w.Code, "GET /v1/games should return HTTP Status 200")
	}

	t.Log("\t And match the created game, with hasNext as false")
	{
		assert.False(t, resp.HasNext, "hasNext should be FALSE")
		assert.Equal(t, 1, len(resp.Games), "returned more than one item")

		returnedGame := resp.Games[0]

		assert.Equal(t, game.Name, returnedGame.Name, "Game Name does not match")
		assert.Equal(t, game.Year, returnedGame.Year, "Game Year does not match")
		assert.Equal(t, game.Platform, returnedGame.Platform, "Game Platform does not match")
		assert.Equal(t, game.Genre, returnedGame.Genre, "Game Genre does not match")
		assert.Equal(t, game.Publisher, returnedGame.Publisher, "Game Publisher does not match")
	}
}
