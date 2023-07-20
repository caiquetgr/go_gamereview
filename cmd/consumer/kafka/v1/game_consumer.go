package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	NewGameTopic = "new-game-event"
)

type NewGameEventHandler struct {
	GameService games.GameService
}

func BuildNewGameEventHandler(gs games.GameService) NewGameEventHandler {
	return NewGameEventHandler{
		GameService: gs,
	}
}

func (ngh NewGameEventHandler) HandleMessage(c *kafka.Consumer, e *kafka.Message) {
	ng := &games.NewGame{}

	if err := json.Unmarshal(e.Value, ng); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling message %v - error %v", e.Value, err)
		return
	}

	game, err := ngh.GameService.CreateGame(context.Background(), *ng)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating game: %v", err)
	} else {
		log.Println("created game", game)
	}

	_, err = c.CommitMessage(e)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error commiting message: %v", err)
	}
}

func (ngh NewGameEventHandler) GetTopic() string {
	return NewGameTopic
}
