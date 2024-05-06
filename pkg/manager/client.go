package manager

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/Hank-Kuo/chat-app/config"
	"github.com/Hank-Kuo/chat-app/pkg/utils"

	"github.com/gorilla/websocket"
)

type Client struct {
	ClientId    string
	Socket      *websocket.Conn
	ConnectTime uint64
}

type SendReq struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type ToClientInfo struct {
	ClientId string
	Data     interface{}
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Socket: conn,
	}
}

func (c *Client) ValidAuth(cfg *config.Config) error {

	messageType, message, err := c.Socket.ReadMessage()
	if err != nil {
		if messageType == -1 && websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
			return err
		}
	}

	msg := string(message)

	if !strings.HasPrefix(msg, "Authorization: Bearer") {
		return errors.New("Invalid format")
	}

	tokenStr := strings.Split(msg, "Bearer ")
	if len(tokenStr) < 1 {
		return errors.New("Invalid format")
	}
	claims, err := utils.ValidJwt(cfg, tokenStr[1])
	if err != nil {
		return err
	}
	c.ClientId = claims.UserID

	return nil
}
