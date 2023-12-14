package message

import (
	"context"

	"github.com/gocql/gocql"

	"chat-app/internal/models"
)

type Repository interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	CreateReply(ctx context.Context, reply *models.Reply) error
	GetMessage(ctx context.Context, channelID string, cursor int64, limit int) ([]*models.Message, error)
	GetReply(ctx context.Context, messageID int64) ([]*models.Reply, error)
}

type messageRepo struct {
	session *gocql.Session
}

func NewRepo(session *gocql.Session) Repository {
	return &messageRepo{session: session}
}
