package test

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ConsumeSingleMessage(c *kafka.Consumer, topic string) ([]byte, error) {
	err := c.SubscribeTopics([]string{topic}, nil)
	defer c.Close()
	if err != nil {
		return nil, err
	}

	ev := c.Poll(1000)

	switch e := ev.(type) {
	case *kafka.Message:
		log.Printf("-- Message on %s: %s\n", e.TopicPartition, string(e.Value))
		log.Printf("-- Headers: %s\n", e.Headers)
		return e.Value, nil
	case kafka.Error:
		return nil, fmt.Errorf("kafka error: %v: %v\n", e.Code(), e)
	default:
		return nil, fmt.Errorf("message not found")
	}
}
