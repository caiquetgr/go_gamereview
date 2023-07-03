package tests

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/caiquetgr/go_gamereview/foundation/test"
	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/caiquetgr/go_gamereview/internal/platform/database"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type GameTest struct {
	it test.AppIntegrationTest
}

func TestGames(t *testing.T) {
	ctx := context.Background()
	it := test.NewIntegrationTest(ctx, comp)

	t.Cleanup(it.Teardown)
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
	topic := "new-game-event"
	kp := g.it.Kp
	ng := games.NewGame{
		Name:      "Super Ghouls'n Ghosts",
		Year:      1991,
		Platform:  "Super Nintendo",
		Genre:     "Platform",
		Publisher: "Capcom",
	}
	bytes, err := json.Marshal(ng)
	if err != nil {
		t.Fatalf("[ERROR] Could not marshal new game to json: %v", err)
	}

	eventCh := make(chan kafka.Event)
	defer close(eventCh)

	kp.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: bytes,
	}, eventCh)

	m := (<-eventCh).(*kafka.Message)

	if m.TopicPartition.Error != nil {
		t.Fatalf("[ERROR] Failed to deliver message: %v", m.TopicPartition.Error.Error())
	}
}
