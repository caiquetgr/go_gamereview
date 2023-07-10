package tests

import (
	"bytes"
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
	"github.com/caiquetgr/go_gamereview/internal/platform/database"
	"github.com/caiquetgr/go_gamereview/internal/platform/events/kafka"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type GameTest struct {
	app http.Handler
	gs  games.GameService
	it  test.AppIntegrationTest
}

type GetGamesV1Response struct {
	Games   []games.Game
	HasNext bool
}

func TestGames(t *testing.T) {
	ctx := context.Background()
	it := test.NewIntegrationTest(ctx)

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
		it: it,
	}

	gameTests := map[string]func(t *testing.T){
		"GetGamesList": tests.TestGetGamesList,
		"CreateGame":   tests.TestCreateGameAsync,
	}

	for k, v := range gameTests {
		tests.BeforeRun()
		t.Run(k, v)
	}
}

func (g GameTest) BeforeRun() {
	g.CleanDatabase()
}

func (g GameTest) CleanDatabase() {
	ctx := context.Background()
	_ = database.Rollback(ctx, g.it.Db)
	_ = database.Migrate(ctx, g.it.Db)
}

func (g GameTest) TestGetGamesList(t *testing.T) {
	game := games.NewGame{
		Name:      "Donkey Kong Country 2: Diddy's Kong Quest",
		Year:      1995,
		Platform:  "Super Nintendo",
		Genre:     "Platform",
		Publisher: "Nintendo",
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

func (g GameTest) TestCreateGameAsync(t *testing.T) {
	game := games.NewGame{
		Name:      "Donkey Kong Country 2: Diddy's Kong Quest",
		Year:      1995,
		Platform:  "Super Nintendo",
		Genre:     "Platform",
		Publisher: "Nintendo",
	}

	rb, err := json.Marshal(game)
	if err != nil {
		t.Errorf("unable to marshal game: %v", err)
	}

	t.Log("\t Given a new game creation request")
	{
		r := httptest.NewRequest(http.MethodPost, "/v1/games", bytes.NewReader(rb))
		w := httptest.NewRecorder()
		g.app.ServeHTTP(w, r)

		assert.Equal(t, http.StatusAccepted, w.Result().StatusCode, "Game POST did not returned 202")
	}

	t.Log("\t Should have posted a new game event equal to game request")
	{
		value, err := kafka.ConsumeSingleMessage(g.it.Kc, "new-game-event")
		if err != nil {
			t.Fatalf("[ERROR] Failed consuming kafka message: %v", err)
		}

		ng := &games.NewGame{}

		if err := json.Unmarshal(value, ng); err != nil {
			t.Fatalf("[ERROR] Failed unmarshalling kafka message: %v", err)
		}

		assert.True(t, cmp.Equal(game, *ng), "New game event is not equal to game request")
	}
}
