package kafka

import (
	"github.com/Hank-Kuo/chat-app/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	Producer  *kafka.Producer
	eventChan chan kafka.Event
}

type KafkaConsumer struct {
	Consumer *kafka.Consumer
}

func NewKafkaProducer(cfg config.KafkaConfig) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.Producer.Brokers,
		"acks":               cfg.Producer.Acks,
		"message.timeout.ms": 5000,
		"enable.idempotence": cfg.Producer.Idepotence,
	})

	if err != nil {
		return nil, err
	}
	return &KafkaProducer{
		Producer:  p,
		eventChan: make(chan kafka.Event, 10000),
	}, nil
}

func (p *KafkaProducer) Produce(message *kafka.Message) error {
	err := p.Producer.Produce(message, p.eventChan)
	if err != nil {
		return err
	}
	e := <-p.eventChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}
	return nil
}

func NewKafkaConsumer(cfg config.KafkaConfig) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        cfg.Consumer.Brokers,
		"broker.address.family":    "v4",
		"group.id":                 cfg.Consumer.GroupID,
		"auto.offset.reset":        cfg.Consumer.OffsetReset,
		"max.poll.interval.ms":     130000,
		"heartbeat.interval.ms":    5000,
		"enable.auto.offset.store": cfg.Consumer.AutoOffset,
	})

	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{Consumer: c}, nil
}

func NewKafkaAdmin(cfg config.KafkaConfig) (*kafka.AdminClient, error) {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Consumer.Brokers,
	})

	if err != nil {
		return nil, err
	}

	return adminClient, nil
}
