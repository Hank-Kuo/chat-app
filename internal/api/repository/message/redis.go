package message

import (
	"context"
	"encoding/json"

	"github.com/Hank-Kuo/chat-app/internal/models"
	"github.com/Hank-Kuo/chat-app/pkg/tracer"

	"github.com/redis/go-redis/v9"
)

func (r *messageRepo) GetChannelsCache(ctx context.Context) ([]*models.Channel, error) {
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

func (r *messageRepo) CreateChannelsCache(ctx context.Context, channels []*models.Channel) error {
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
