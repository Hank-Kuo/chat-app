package channel

import (
	"context"

	"github.com/Hank-Kuo/chat-app/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	Get(ctx context.Context) ([]*models.Channel, error)
	GetByID(ctx context.Context, channelID string) (*models.Channel, error)
	Create(ctx context.Context, channel *models.Channel) (string, error)
	CreateUserToChannel(ctx context.Context, uchannel *models.UserToChannel) error
	GetUserToChannel(ctx context.Context, userID string) ([]*models.Channel, error)
	GetUserByChannel(ctx context.Context, channelID string) ([]*models.UserToChannel, error)

	GetChannelsCache(ctx context.Context) ([]*models.Channel, error)
	CreateChannelsCache(ctx context.Context, channels []*models.Channel) error

	GetUserByChannelCache(ctx context.Context, channelID string) ([]*models.UserToChannel, error)
	CreateUserByChannelCache(ctx context.Context, channelID string, uchannel []*models.UserToChannel) error

	GetChannelsByUserCache(ctx context.Context, userID string) ([]*models.Channel, error)
	CreateChannelsByUserCache(ctx context.Context, userID string, channels []*models.Channel) error
}

type channelRepo struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewRepo(db *sqlx.DB, rdb *redis.Client) Repository {
	return &channelRepo{db: db, rdb: rdb}
}
