package tests

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	kafkahandler "github.com/caiquetgr/go_gamereview/cmd/consumer/kafka"
	"github.com/caiquetgr/go_gamereview/foundation/test"
	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/caiquetgr/go_gamereview/internal/domain/games/db/gamedb"
	"github.com/caiquetgr/go_gamereview/internal/platform/database"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"
)

type GameTest struct {
	it test.AppIntegrationTest
	r  games.Repository
}

func TestGames(t *testing.T) {
	ctx := context.Background()
	it := test.NewIntegrationTest(ctx, comp)
	t.Cleanup(it.Teardown)

	gt := GameTest{
		it: it,
		r:  gamedb.NewGameRepositoryBun(it.Db),
	}

	tests := map[string]func(*testing.T){
		"TestCreateGame": gt.TestCreateGame,
	}

	sigChan := make(chan os.Signal, 1)
	stopChan := make(chan struct{})
	defer func() {
		close(sigChan)
	}()

	go kafkahandler.Handle(kafkahandler.KafkaHandlerConfig{
		DB:                  it.Db,
		KafkaProducer:       it.Kp,
		KafkaConsumerCreate: it.KcCreator,
		SigChan:             sigChan,
		StopChan:            stopChan,
	})

	for k, v := range tests {
		gt.BeforeRun()
		t.Run(k, v)
	}

	close(stopChan)
}

func (g GameTest) BeforeRun() {
	g.CleanDatabase()
}

func (g GameTest) CleanDatabase() {
	ctx := context.Background()
	_ = database.Rollback(ctx, g.it.Db)
	_ = database.Migrate(ctx, g.it.Db)
}

func (g GameTest) TestCreateGame(t *testing.T) {
	ctx := context.Background()
	topic := "new-game-event"
	kp := g.it.Kp
	ng := games.NewGame{
		Name:      "Super Ghouls'n Ghosts",
		Year:      1991,
		Platform:  "Super Nintendo",
		Genre:     "Platform",
		Publisher: "Capcom",
	}

	t.Log("Given a new game event")
	{
		bytes, err := json.Marshal(ng)
		if err != nil {
			t.Fatalf("[ERROR] Could not marshal new game to json: %v", err)
		}

		eventCh := make(chan kafka.Event)
		defer close(eventCh)

		err = kp.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: bytes,
		}, eventCh)

		if err != nil {
			t.Fatalf("[ERROR] Failed to produce event: %v", err)
		}

		m := (<-eventCh).(*kafka.Message)

		if m.TopicPartition.Error != nil {
			t.Fatalf("[ERROR] Failed to deliver message: %v", m.TopicPartition.Error.Error())
		}
	}

	t.Log("Should create a game successfully")
	{
		f := func() bool {
			game, err := g.r.FindByName(ctx, ng.Name)
			if err != nil {
				return false
			}
			return game.Name != ""
		}

		success, err := test.WaitUntil(f, 5*time.Second)

		if !success || err != nil {
			t.Fatalf("[ERROR] Couldn't find the game to continue the test: success=%v, err=%v", success, err)
		}
	}

	t.Log("And match the game event")
	{
		game, _ := g.r.FindByName(ctx, "Super Ghouls'n Ghosts")
		assert.Equal(t, ng.Name, game.Name, "Game Name does not match")
		assert.Equal(t, ng.Year, game.Year, "Game Year does not match")
		assert.Equal(t, ng.Platform, game.Platform, "Game Platform does not match")
		assert.Equal(t, ng.Genre, game.Genre, "Game Genre does not match")
		assert.Equal(t, ng.Publisher, game.Publisher, "Game Publisher does not match")
	}
}
