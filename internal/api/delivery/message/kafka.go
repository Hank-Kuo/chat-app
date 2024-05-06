package message

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Hank-Kuo/chat-app/config"
	channelSrv "github.com/Hank-Kuo/chat-app/internal/api/service/channel"
	messageSrv "github.com/Hank-Kuo/chat-app/internal/api/service/message"
	"github.com/Hank-Kuo/chat-app/internal/models"
	"github.com/Hank-Kuo/chat-app/pkg/kafka"
	"github.com/Hank-Kuo/chat-app/pkg/logger"
	"github.com/Hank-Kuo/chat-app/pkg/manager"

	kafka_go "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type kafkaHandler struct {
	ctx              context.Context
	cfg              *config.Config
	messageSrv       messageSrv.Service
	channelSrv       channelSrv.Service
	manager          *manager.ClientManager
	subscribeChannel []string
	logger           logger.Logger
}

func NewKafkaHandler(cfg *config.Config, messageSrv messageSrv.Service, channelSrv channelSrv.Service, manager *manager.ClientManager, logger logger.Logger) *kafkaHandler {
	handler := &kafkaHandler{
		ctx:              context.Background(),
		cfg:              cfg,
		messageSrv:       messageSrv,
		channelSrv:       channelSrv,
		manager:          manager,
		subscribeChannel: []string{},
		logger:           logger,
	}
	return handler
}

// update subscriber channel_topic
func (h *kafkaHandler) Listen(c *kafka.KafkaConsumer, admin *kafka_go.AdminClient) {
	for {
		select {
		case <-time.After(1 * time.Second):
			if len(h.manager.ClientIdMap) > 0 {
				meta, err := admin.GetMetadata(nil, true, 2000)
				if err != nil {
					h.logger.Error(err)
					continue
				}
				totalChannelsMap := map[string]bool{}

				for u, _ := range h.manager.ClientIdMap {
					if channels, err := h.channelSrv.GetChannelByUser(h.ctx, u); err == nil {
						for _, channel := range channels {
							totalChannelsMap["message_"+channel.ID] = true
						}
					}
				}

				totalChannels := []string{}
				for _, t := range meta.Topics {
					if _, ok := totalChannelsMap[t.Topic]; ok {
						totalChannels = append(totalChannels, t.Topic)
					}
				}

				if len(totalChannels) > 0 {
					if len(h.subscribeChannel) != len(totalChannels) {
						h.subscribeChannel = totalChannels
						if err := c.Consumer.SubscribeTopics(totalChannels, nil); err != nil {
							h.logger.Error(err)
						}
					}
				}
			}
		}
	}
}

// given channel, publish message to user
func (h *kafkaHandler) ConsumeMessage(c *kafka.KafkaConsumer) {
	for {
		msg, err := c.Consumer.ReadMessage(100)
		if err != nil {
			continue
		}

		switch string(msg.Key) {
		case "topic_message":
			if err := h.consumeMessage(string(msg.Value)); err != nil {
				h.logger.Errorf("broadcast message got error: %v", err)
			}
			_, err = c.Consumer.StoreMessage(msg)
			if err != nil {
				h.logger.Errorf("broadcast reply got error: %v", err)

			}
		case "topic_reply":
			if err := h.consumeReply(string(msg.Value)); err != nil {
				h.logger.Errorf("broadcast reply got error: %v", err)
			}
			_, err = c.Consumer.StoreMessage(msg)
			if err != nil {
				h.logger.Errorf("broadcast reply got error: %v", err)

			}
		}
	}
}

func (h *kafkaHandler) consumeMessage(data string) error {
	var body models.Message
	if err := json.Unmarshal([]byte(data), &body); err != nil {
		return err
	}

	users, err := h.channelSrv.GetUserByChannel(h.ctx, body.ChannelID)
	if err != nil {
		return err
	}
	for _, u := range users {
		if u.UserID != body.UserID {
			h.manager.ToClientChan <- manager.ToClientInfo{ClientId: u.UserID, Data: body}
		}
	}
	return nil
}

func (h *kafkaHandler) consumeReply(data string) error {
	var body models.Reply
	if err := json.Unmarshal([]byte(data), &body); err != nil {
		return err
	}

	users, err := h.channelSrv.GetUserByChannel(h.ctx, body.ChannelID)
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.UserID != body.UserID {
			h.manager.ToClientChan <- manager.ToClientInfo{ClientId: u.UserID, Data: body}
		}
	}
	return nil
}
