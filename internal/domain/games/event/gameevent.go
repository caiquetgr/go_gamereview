package event

import (
	"context"
	"encoding/json"
	"log"

	"github.com/caiquetgr/go_gamereview/internal/domain/games"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type GameEventProducer struct {
	topic string
	p     *kafka.Producer
}

func NewGameEventProducer(topic string, p *kafka.Producer) GameEventProducer {
	return GameEventProducer{
		topic: topic,
		p:     p,
	}
}

func (ge GameEventProducer) CreateGameEvent(ctx context.Context, ng games.NewGame) error {
	bytes, err := json.Marshal(ng)

	if err != nil {
		return err
	}

	eventCh := make(chan kafka.Event)
	defer close(eventCh)

	err = ge.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &ge.topic,
			Partition: kafka.PartitionAny,
		},
		Value: bytes,
	}, eventCh)

	if err != nil {
		return err
	}

	// TODO: what if kafka is down?
	m := (<-eventCh).(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("failed to deliver message: %v\n", m.TopicPartition.Error.Error())
		return m.TopicPartition.Error
	} else {
		log.Printf("delivered message to topic %s [%d] offset %v: %+v\n", *m.TopicPartition.Topic,
			m.TopicPartition.Partition, m.TopicPartition.Offset, ng)
	}

	return nil
}
