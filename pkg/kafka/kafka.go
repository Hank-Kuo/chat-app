package kafka

import (
	"errors"

	"github.com/Hank-Kuo/chat-app/config"

	"github.com/segmentio/kafka-go"
)

func NewWriter(cfg config.KafkaConfig, topicName string) (*kafka.Writer, error) {
	topic, ok := cfg.Topics[topicName]

	if !ok {
		return nil, errors.New("not found topic")
	}

	return &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    topic.Name,
		Balancer: &kafka.LeastBytes{},
	}, nil
}
