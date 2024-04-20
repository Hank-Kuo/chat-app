package channel

import (
	"chat-app/internal/models"
	"chat-app/pkg/tracer"
	"context"

	"github.com/pkg/errors"
)

func (r *channelRepo) Get(ctx context.Context) ([]*models.Channel, error) {
	ctx, span := tracer.NewSpan(ctx, "ChannelRepo.Get", nil)
	defer span.End()

	channels := []*models.Channel{}
	if err := r.db.SelectContext(ctx, &channels, "SELECT * FROM channel"); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "ChannelRepo.Get")
	}

	return channels, nil
}

func (r *channelRepo) Create(ctx context.Context, channel *models.Channel) (string, error) {
	ctx, span := tracer.NewSpan(ctx, "ChannelRepo.Create", nil)
	defer span.End()

	var channelID string
	sqlQuery := `INSERT INTO channel(name, user_id) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRowxContext(ctx, sqlQuery, channel.Name, channel.UserID).Scan(&channelID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return "", errors.Wrap(err, "ChannelRepo.Create")
	}
	return channelID, nil
}

func (r *channelRepo) CreateUserToChannel(ctx context.Context, uchannel *models.UserToChannel) error {
	ctx, span := tracer.NewSpan(ctx, "ChannelRepo.CreateUserToChannel", nil)
	defer span.End()

	sqlQuery := `INSERT INTO user_to_channel(channel_id, user_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, sqlQuery, uchannel.ChannelID, uchannel.UserID)

	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "ChannelRepo.CreateUserToChannel")
	}
	return nil
}
func (r *channelRepo) GetUserToChannel(ctx context.Context, userID string) ([]*models.Channel, error) {
	ctx, span := tracer.NewSpan(ctx, "ChannelRepo.GetUserToChannel", nil)
	defer span.End()

	sqlQuery := `
		SELECT channel.*
		FROM user_to_channel
		INNER JOIN channel ON channel.id = user_to_channel.channel_id
		WHERE user_to_channel.user_id = $1
	`
	channels := []*models.Channel{}
	if err := r.db.SelectContext(ctx, &channels, sqlQuery, userID); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "ChannelRepo.GetUserToChannel")
	}

	return channels, nil

}

func (r *channelRepo) GetUserByChannel(ctx context.Context, channelID string) ([]*models.UserToChannel, error) {
	ctx, span := tracer.NewSpan(ctx, "ChannelRepo.GetUserByChannel", nil)
	defer span.End()

	sqlQuery := `
		SELECT *
		FROM user_to_channel
		WHERE channel_id = $1
	`
	channels := []*models.UserToChannel{}
	if err := r.db.SelectContext(ctx, &channels, sqlQuery, channelID); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "ChannelRepo.GetUserByChannel")
	}

	return channels, nil

}
