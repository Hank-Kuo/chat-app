package channel

import (
	"context"
	"encoding/json"

	"github.com/Hank-Kuo/chat-app/internal/models"
	"github.com/Hank-Kuo/chat-app/pkg/tracer"

	"github.com/redis/go-redis/v9"
)

func (r *channelRepo) CreateChannelsByUserCache(ctx context.Context, userID string, channels []*models.Channel) error {
	ctx, span := tracer.NewSpan(ctx, "ChannelRepo.CreateUsersChannelsCache", nil)
	defer span.End()

	channelsByte, err := r.rdb.Get(ctx, userID).Result()

	oldChannels := []*models.Channel{}

	if err == redis.Nil {
		oldChannels = append(oldChannels, channels...)
	} else if err != nil {
		tracer.AddSpanError(span, err)
		return err
	} else {
		if err = json.Unmarshal([]byte(channelsByte), &oldChannels); err != nil {
			tracer.AddSpanError(span, err)
			return err
		}
		oldChannels = append(oldChannels, channels...)
	}

	jsonData, err := json.Marshal(oldChannels)
	if err != nil {
		tracer.AddSpanError(span, err)
		return err
	}
	if err := r.rdb.Set(ctx, userID, jsonData, 0).Err(); err != nil {
		tracer.AddSpanError(span, err)
		return err
	}

	return nil
}

func (r *channelRepo) CreateUserByChannelCache(ctx context.Context, channelID string, uchannel []*models.UserToChannel) error {
	ctx, span := tracer.NewSpan(ctx, "ChannelRepo.CreateUserByChannelCache", nil)
	defer span.End()

	userToChannelByte, err := r.rdb.Get(ctx, channelID).Result()

	userToChannel := []*models.UserToChannel{}

	if err == redis.Nil {
		userToChannel = append(userToChannel, uchannel...)
	} else if err != nil {
		tracer.AddSpanError(span, err)
		return err
	} else {
		if err = json.Unmarshal([]byte(userToChannelByte), &userToChannel); err != nil {
			tracer.AddSpanError(span, err)
			return err
		}
		userToChannel = append(userToChannel, uchannel...)
	}

	jsonData, err := json.Marshal(userToChannel)
	if err != nil {
		tracer.AddSpanError(span, err)
		return err
	}
	if err = r.rdb.Set(ctx, channelID, jsonData, 0).Err(); err != nil {
		tracer.AddSpanError(span, err)
		return err
	}

	return nil
}

func (r *channelRepo) GetChannelsByUserCache(ctx context.Context, userID string) ([]*models.Channel, error) {
	ctx, span := tracer.NewSpan(ctx, "ChannelRepo.GetChannelsByUserCache", nil)
	defer span.End()

	channels := []*models.Channel{}

	channelsByte, err := r.rdb.Get(ctx, userID).Result()
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	if err = json.Unmarshal([]byte(channelsByte), &channels); err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return channels, nil
}

func (r *channelRepo) GetUserByChannelCache(ctx context.Context, channelID string) ([]*models.UserToChannel, error) {
	ctx, span := tracer.NewSpan(ctx, "ChannelRepo.GetUserByChannelCache", nil)
	defer span.End()

	usersToChannel := []*models.UserToChannel{}

	usersToChannelByte, err := r.rdb.Get(ctx, channelID).Result()
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	if err = json.Unmarshal([]byte(usersToChannelByte), &usersToChannel); err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return usersToChannel, nil
}

func (r *channelRepo) GetChannelsCache(ctx context.Context) ([]*models.Channel, error) {
	ctx, span := tracer.NewSpan(ctx, "ChannelRepo.GetChannelCache", nil)
	defer span.End()

	channels := []*models.Channel{}

	channelsByte, err := r.rdb.Get(ctx, "channels").Result()
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	if err = json.Unmarshal([]byte(channelsByte), &channels); err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return channels, nil
}

func (r *channelRepo) CreateChannelsCache(ctx context.Context, channels []*models.Channel) error {
	ctx, span := tracer.NewSpan(ctx, "ChannelRepo.CreateChannelsCache", nil)
	defer span.End()

	channelsByte, err := r.rdb.Get(ctx, "channels").Result()

	oldChannels := []*models.Channel{}

	if err == redis.Nil {
		oldChannels = append(oldChannels, channels...)
	} else if err != nil {
		tracer.AddSpanError(span, err)
		return err
	} else {
		if err = json.Unmarshal([]byte(channelsByte), &oldChannels); err != nil {
			tracer.AddSpanError(span, err)
			return err
		}
		oldChannels = append(oldChannels, channels...)
	}

	jsonData, err := json.Marshal(oldChannels)
	if err != nil {
		tracer.AddSpanError(span, err)
		return err
	}

	if err := r.rdb.Set(ctx, "channels", jsonData, 0).Err(); err != nil {
		tracer.AddSpanError(span, err)
		return err
	}

	return nil
}
