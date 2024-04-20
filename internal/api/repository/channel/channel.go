package channel

import (
	"context"

	"github.com/Hank-Kuo/chat-app/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Get(ctx context.Context) ([]*models.Channel, error)
	Create(ctx context.Context, channel *models.Channel) (string, error)
	CreateUserToChannel(ctx context.Context, uchannel *models.UserToChannel) error
	GetUserToChannel(ctx context.Context, userID string) ([]*models.Channel, error)
	GetUserByChannel(ctx context.Context, channelID string) ([]*models.UserToChannel, error)
}

type channelRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &channelRepo{db: db}
}
