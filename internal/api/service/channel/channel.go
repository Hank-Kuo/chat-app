package channel

import (
	"context"

	"chat-app/config"
	channelRepo "chat-app/internal/api/repository/channel"
	"chat-app/internal/models"
	// "chat-app/pkg/customError"
	"chat-app/pkg/logger"
	"chat-app/pkg/tracer"

	"github.com/pkg/errors"
)

type Service interface {
	Get(ctx context.Context) ([]*models.Channel, error)
	Create(ctx context.Context, channel *models.Channel) error
	Join(ctx context.Context, userID string, channelID string) error
	GetUserChannel(ctx context.Context, userID string) ([]*models.Channel, error)
}

type channelSrv struct {
	cfg         *config.Config
	channelRepo channelRepo.Repository
	logger      logger.Logger
}

func NewService(cfg *config.Config, channelRepo channelRepo.Repository, logger logger.Logger) Service {
	return &channelSrv{
		cfg:         cfg,
		channelRepo: channelRepo,
		logger:      logger,
	}
}

func (srv *channelSrv) Get(ctx context.Context) ([]*models.Channel, error) {
	c, span := tracer.NewSpan(ctx, "ChannelService.Get", nil)
	defer span.End()

	channels, err := srv.channelRepo.Get(c)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "ChannelService.Get")
	}
	return channels, nil
}

func (srv *channelSrv) Create(ctx context.Context, channel *models.Channel) error {
	c, span := tracer.NewSpan(ctx, "ChannelService.Create", nil)
	defer span.End()

	channelID, err := srv.channelRepo.Create(c, channel)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "ChannelService.Create")
	}

	if err := srv.channelRepo.CreateUserToChannel(c,
		&models.UserToChannel{
			UserID: channel.UserID, ChannelID: channelID,
		}); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "ChannelService.Create")
	}
	return nil
}
func (srv *channelSrv) Join(ctx context.Context, userID, channelID string) error {
	c, span := tracer.NewSpan(ctx, "ChannelService.Join", nil)
	defer span.End()

	if err := srv.channelRepo.CreateUserToChannel(c,
		&models.UserToChannel{
			UserID: userID, ChannelID: channelID,
		}); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "ChannelService.Join")
	}

	return nil
}

func (srv *channelSrv) GetUserChannel(ctx context.Context, userID string) ([]*models.Channel, error) {
	c, span := tracer.NewSpan(ctx, "ChannelService.GetUserChannel", nil)
	defer span.End()

	channels, err := srv.channelRepo.GetUserToChannel(c, userID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "ChannelService.GetUserChannel")
	}
	return channels, nil
}
