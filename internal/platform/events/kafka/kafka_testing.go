package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func ConsumeSingleMessage(c *kafka.Consumer, topic string) ([]byte, error) {
	err := c.SubscribeTopics([]string{topic}, nil)
	defer c.Unsubscribe()

	if err != nil {
		return nil, err
	}

	run := true
	var v []byte

	for run {
		ev := c.Poll(100)

		select {
		case <-time.After(5 * time.Second):
			v, err = nil, fmt.Errorf("Timeout consume single message")
			run = false
		default:
			switch e := ev.(type) {
			case *kafka.Message:
				log.Printf("-- Message on %s: %s\n", e.TopicPartition, string(e.Value))
				log.Printf("-- Headers: %s\n", e.Headers)
				v, err = e.Value, nil
				run = false
			case kafka.Error:
				v, err = nil, fmt.Errorf("kafka error: %v: %v\n", e.Code(), e)
				run = false
			default:
				v, err = nil, fmt.Errorf("message not found")
			}
		}
	}

	return v, err
}
