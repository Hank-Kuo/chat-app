package message

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/redis/go-redis/v9"

	"github.com/Hank-Kuo/chat-app/internal/models"
	"github.com/Hank-Kuo/chat-app/pkg/kafka"
)

type Repository interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	CreateReply(ctx context.Context, reply *models.Reply) error
	GetMessage(ctx context.Context, channelID string, cursor int64, limit int) ([]*models.Message, error)
	GetReply(ctx context.Context, messageID int64, cursor int64, limit int) ([]*models.Reply, error)
	PublishMessage(ctx context.Context, message *models.Message) error
	PublishReply(ctx context.Context, reply *models.Reply) error
}

type messageRepo struct {
	session       *gocql.Session
	kafkaProducer *kafka.KafkaProducer
	rdb           *redis.Client
}

func NewRepo(session *gocql.Session, kafkaProducer *kafka.KafkaProducer, rdb *redis.Client) Repository {
	return &messageRepo{session: session, kafkaProducer: kafkaProducer, rdb: rdb}
}
