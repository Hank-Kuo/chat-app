package message

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/segmentio/kafka-go"

	"github.com/Hank-Kuo/chat-app/internal/models"
)

type Repository interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	CreateReply(ctx context.Context, reply *models.Reply) error
	GetMessage(ctx context.Context, channelID string, cursor int64, limit int) ([]*models.Message, error)
	GetReply(ctx context.Context, messageID int64, cursor int64, limit int) ([]*models.Reply, error)
	PublishMessage(ctx context.Context, userID string, message *models.Message) error
	PublishReply(ctx context.Context, userID string, reply *models.Reply) error
}

type messageRepo struct {
	session            *gocql.Session
	kafkaMessageWriter *kafka.Writer
	kafkaReplyWriter   *kafka.Writer
}

func NewRepo(session *gocql.Session, kafkaMessageWriter *kafka.Writer, kafkaReplyWriter *kafka.Writer) Repository {
	return &messageRepo{session: session, kafkaMessageWriter: kafkaMessageWriter, kafkaReplyWriter: kafkaReplyWriter}
}
