package channel

import (
	"context"

	"github.com/Hank-Kuo/chat-app/config"
	channelRepo "github.com/Hank-Kuo/chat-app/internal/api/repository/channel"
	"github.com/Hank-Kuo/chat-app/internal/models"
	"github.com/Hank-Kuo/chat-app/pkg/logger"
	"github.com/Hank-Kuo/chat-app/pkg/tracer"

	"github.com/pkg/errors"
)

type Service interface {
	Get(ctx context.Context) ([]*models.Channel, error)
	Create(ctx context.Context, channel *models.Channel) error
	Join(ctx context.Context, userID string, channelID string) error
	GetChannelByUser(ctx context.Context, userID string) ([]*models.Channel, error)
	GetUserByChannel(ctx context.Context, channelID string) ([]*models.UserToChannel, error)
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

	if channels, err := srv.channelRepo.GetChannelsCache(c); err == nil {
		return channels, nil
	}

	channels, err := srv.channelRepo.Get(c)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "ChannelService.Get")
	}
	if err = srv.channelRepo.CreateChannelsCache(c, channels); err != nil {
		tracer.AddSpanError(span, err)
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

	channel.ID = channelID
	if err = srv.channelRepo.CreateChannelsCache(c, []*models.Channel{channel}); err != nil {
		tracer.AddSpanError(span, err)
	}

	uchannel := &models.UserToChannel{
		UserID: channel.UserID, ChannelID: channelID,
	}
	if err := srv.channelRepo.CreateUserToChannel(c, uchannel); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "ChannelService.Create")
	}

	if err = srv.channelRepo.CreateUserByChannelCache(c, channelID, []*models.UserToChannel{uchannel}); err != nil {
		tracer.AddSpanError(span, err)
	}
	if err := srv.channelRepo.CreateChannelsByUserCache(c, channel.UserID, []*models.Channel{channel}); err != nil {
		tracer.AddSpanError(span, err)
	}

	return nil
}

func (srv *channelSrv) Join(ctx context.Context, userID, channelID string) error {
	c, span := tracer.NewSpan(ctx, "ChannelService.Join", nil)
	defer span.End()

	uchannel := &models.UserToChannel{
		UserID: userID, ChannelID: channelID,
	}

	if err := srv.channelRepo.CreateUserByChannelCache(c, channelID, []*models.UserToChannel{uchannel}); err != nil {
		tracer.AddSpanError(span, err)
	}

	if err := srv.channelRepo.CreateUserToChannel(c, uchannel); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "ChannelService.Join")
	}

	channel, err := srv.channelRepo.GetByID(c, channelID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "ChannelService.Join")
	}

	if err := srv.channelRepo.CreateChannelsByUserCache(c, userID, []*models.Channel{channel}); err != nil {
		tracer.AddSpanError(span, err)
	}

	return nil
}

func (srv *channelSrv) GetChannelByUser(ctx context.Context, userID string) ([]*models.Channel, error) {
	c, span := tracer.NewSpan(ctx, "ChannelService.GetChannelByUser", nil)
	defer span.End()

	channels, err := srv.channelRepo.GetUserToChannel(c, userID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "ChannelService.GetChannelByUser")
	}

	if err := srv.channelRepo.CreateChannelsByUserCache(c, userID, channels); err != nil {
		tracer.AddSpanError(span, err)
	}
	return channels, nil
}

func (srv *channelSrv) GetUserByChannel(ctx context.Context, channelID string) ([]*models.UserToChannel, error) {
	c, span := tracer.NewSpan(ctx, "ChannelService.GetUserByChannel", nil)
	defer span.End()
	if uchannel, err := srv.channelRepo.GetUserByChannelCache(c, channelID); err == nil {
		return uchannel, nil
	}
	uchannel, err := srv.channelRepo.GetUserByChannel(c, channelID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "ChannelService.GetUserByChannel")
	}
	if err := srv.channelRepo.CreateUserByChannelCache(c, channelID, uchannel); err != nil {
		tracer.AddSpanError(span, err)
	}

	return uchannel, nil
}
