package kafka

import (
	"errors"

	"github.com/Hank-Kuo/chat-app/config"

	"github.com/segmentio/kafka-go"
)

func NewWriter(cfg config.KafkaConfig, topicName string) (*kafka.Writer, error) {
	for _, t := range cfg.Topics {
		if t.Name == topicName {
			return &kafka.Writer{
				Addr:     kafka.TCP(cfg.Brokers...),
				Topic:    t.Name,
				Balancer: &kafka.LeastBytes{},
			}, nil
		}
	}

	return nil, errors.New("not found topic")

}
