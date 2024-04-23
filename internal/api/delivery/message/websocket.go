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
			tracer.AddSpanError(span, err)
			if messageType == -1 && websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
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
					channels, err := h.channelSrv.GetUserByChannel(ctx, body.ChannelID)

					if err != nil {
						tracer.AddSpanError(span, err)
						httpResponse.Fail(err, h.logger).ToWebSocketJSON(clientSocket.Socket)
					} else {

						// get group client ip
						for _, ch := range channels {
							// if ch.UserID != body.UserID {
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
							// }
						}
						// httpResponse.OK(http.StatusOK, "create message successfully", newMessage).ToWebSocketJSON(clientSocket.Socket)
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
				fmt.Println(clientInfo.InstanceId)
				// send to other instance

			}
		}
	}
}
