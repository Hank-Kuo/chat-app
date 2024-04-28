package message

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hank-Kuo/chat-app/config"
	channelSrv "github.com/Hank-Kuo/chat-app/internal/api/service/channel"
	messageSrv "github.com/Hank-Kuo/chat-app/internal/api/service/message"
	"github.com/Hank-Kuo/chat-app/internal/dto"
	httpMiddleware "github.com/Hank-Kuo/chat-app/internal/middleware/http"
	"github.com/Hank-Kuo/chat-app/pkg/customError"
	"github.com/Hank-Kuo/chat-app/pkg/logger"
	"github.com/Hank-Kuo/chat-app/pkg/manager"
	httpResponse "github.com/Hank-Kuo/chat-app/pkg/response/http_response"
	"github.com/Hank-Kuo/chat-app/pkg/tracer"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type websocketHandler struct {
	cfg        *config.Config
	messageSrv messageSrv.Service
	channelSrv channelSrv.Service
	manager    *manager.ClientManager
	logger     logger.Logger
}

func NewWebSocketHandler(e *gin.RouterGroup, cfg *config.Config, messageSrv messageSrv.Service, channelSrv channelSrv.Service, manager *manager.ClientManager, mid *httpMiddleware.Middleware, logger logger.Logger) {
	handler := &websocketHandler{
		cfg:        cfg,
		messageSrv: messageSrv,
		channelSrv: channelSrv,
		manager:    manager,
		logger:     logger,
	}

	e.GET("/message/ws", handler.Send)
}

func (h *websocketHandler) Send(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "MessageWebsocketHandler.Send", nil)
	defer span.End()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	defer conn.Close()

	if err != nil {
		tracer.AddSpanError(span, err)
		httpResponse.Fail(err, h.logger).ToWebSocketJSON(conn)
		return
	}

	// Establish client and add client into manager
	clientSocket := manager.NewClient(conn)

	if err := clientSocket.ValidAuth(h.cfg); err != nil {
		tracer.AddSpanError(span, err)
		httpResponse.Fail(err, h.logger).ToWebSocketJSON(clientSocket.Socket)
		return
	}

	h.manager.AddClient(clientSocket)
	h.manager.Connect <- clientSocket
	go BroadcastMessage(h.manager)

	for {
		messageType, message, err := clientSocket.Socket.ReadMessage()
		if err != nil {
			if messageType == -1 && websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				tracer.AddSpanError(span, err)
				h.manager.DisConnect <- clientSocket
				return
			}
		}

		var m manager.SendReq

		if err = json.Unmarshal(message, &m); err != nil {
			httpResponse.Fail(err, h.logger).ToWebSocketJSON(clientSocket.Socket)
		} else {
			if m.Action == "CreateMessage" {
				var body dto.CreateMessageReqDto
				if err := json.Unmarshal([]byte(m.Data), &body); err != nil {
					tracer.AddSpanError(span, err)
					httpResponse.Fail(customError.ErrBadRequest, h.logger).ToWebSocketJSON(clientSocket.Socket)
				} else {
					newMessage, err := h.messageSrv.CreateMessage(ctx, &body)

					if err != nil {
						tracer.AddSpanError(span, err)
						httpResponse.Fail(err, h.logger).ToWebSocketJSON(clientSocket.Socket)
					} else {

						httpResponse.OK(http.StatusOK, "create message successfully", newMessage).ToWebSocketJSON(clientSocket.Socket)

						if channels, err := h.channelSrv.GetUserByChannel(ctx, body.ChannelID); err == nil {
							// get group client ip
							for _, ch := range channels {
								if ch.UserID != body.UserID {
									instance, err := h.manager.GetInstacesByClients(ch.UserID)
									if err != nil {
										// offline user
										if err = h.messageSrv.MessageNotification(ctx, ch.UserID, newMessage); err != nil {
											tracer.AddSpanError(span, err)
										}
									} else {
										// online user
										h.manager.ToClientChan <- manager.ToClientInfo{OriginClientId: body.UserID, ClientId: ch.UserID, InstanceId: instance, Data: newMessage}
									}
								}
							}
						}

					}
				}
			} else if m.Action == "CreateReply" {
				var body dto.CreateReplyReqDto
				if err := json.Unmarshal([]byte(m.Data), &body); err != nil {
					tracer.AddSpanError(span, err)
					httpResponse.Fail(customError.ErrBadRequest, h.logger).ToWebSocketJSON(clientSocket.Socket)
				} else {

					newReply, err := h.messageSrv.CreateReply(ctx, &body)

					if err != nil {
						tracer.AddSpanError(span, err)
						httpResponse.Fail(err, h.logger).ToWebSocketJSON(clientSocket.Socket)
					} else {
						httpResponse.OK(http.StatusOK, "create reply successfully", newReply).ToWebSocketJSON(clientSocket.Socket)

						if channels, err := h.channelSrv.GetUserByChannel(ctx, body.ChannelID); err == nil {
							for _, ch := range channels {
								if ch.UserID != body.UserID {
									instance, err := h.manager.GetInstacesByClients(ch.UserID)
									if err != nil {
										// offline user˝
										if err = h.messageSrv.ReplyNotification(ctx, ch.UserID, newReply); err != nil {
											tracer.AddSpanError(span, err)
										}
									} else {
										// online user
										h.manager.ToClientChan <- manager.ToClientInfo{OriginClientId: body.UserID, ClientId: ch.UserID, InstanceId: instance, Data: newReply}
									}
								}

							}
						}

					}
				}
			} else {
				tracer.AddSpanError(span, err)
				httpResponse.Fail(customError.ErrBadRequest, h.logger).ToWebSocketJSON(clientSocket.Socket)
			}
		}
	}

}

func BroadcastMessage(m *manager.ClientManager) {
	for {
		select {
		case clientInfo := <-m.ToClientChan:
			if m.InstanceId == clientInfo.InstanceId {
				// send to local
				if conn, err := m.GetByClientId(clientInfo.ClientId); err == nil && conn != nil {
					httpResponse.OK(http.StatusOK, "send message successfully", clientInfo.Data).ToWebSocketJSON(conn.Socket)
				}
			} else {
				/*
					queryHost := fmt.Sprintf(`%s:%s`, clientInfo.InstanceId, m.cfg.Server.GrpcPort)
					if conn, err := grpc.Dial(queryHost, grpc.WithInsecure()); err == nil {

					}

					c := messagePb.NewMessageServiceClient(conn)

					if _, err = c.MessageReceived(context.Background(), &messagePb.ReceiveMessageRequest{
						OriginClientID: "bac09a89-df1a-4644-ba2f-89f4da8d0456",
						ClientId:       "257e4caf-fb4b-43a1-a4b3-cca94f583bd5",
						InstanceId:     "127.0.0.1",
						Message: &messagePb.ReceiveMessageRequest_Message{
							ChannelId: "b37c4896-70e5-4a94-bbab-7de13e5e41ff",
							MessageId: 1984282932957937700,
							Content:   "send from other",
							UserId:    "bac09a89-df1a-4644-ba2f-89f4da8d0456",
							Username:  "other",
							CreatedAt: "2023-12-01T05:52:17.581988Z",
						},
					}); err != nil {

					}
				*/
				// send to other instance
				fmt.Println(clientInfo.InstanceId)
			}
		}
	}
}
