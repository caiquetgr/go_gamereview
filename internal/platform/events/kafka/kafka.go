package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProducerConfig struct {
	BootstrapServers string
	Acks             string
}

type ConsumerConfig struct {
	BootstrapServers string
	GroupId          string
	AutoOffsetReset  string
}

func CreateKafkaProducer(kpc ProducerConfig) *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kpc.BootstrapServers,
		"acks":              kpc.Acks,
	})
	if err != nil {
		panic(err)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("error delivering event in topic %v\n", ev.TopicPartition)
				} else {
					log.Printf("delivered event in topic %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return p
}

func CreateKafkaConsumer(kcc ConsumerConfig) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kcc.BootstrapServers,
		"group.id":          kcc.GroupId,
		"auto.offset.reset": kcc.AutoOffsetReset,
	})
	if err != nil {
		panic(err)
	}

	return c
}
