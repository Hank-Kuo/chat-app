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
	IsDeleted   bool
}

type SendReq struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type ToClientInfo struct {
	OriginClientId string
	ClientId       string
	InstanceId     string
	Data           interface{}
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

// var ToClientChan = make(chan clientInfo, 1000)

/*
func (c *Client) Read() {
	for {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil {
			if messageType == -1 && websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				// Manager.DisConnect <- c
				return
			} else if messageType != websocket.PingMessage {
				return
			}
		}
		fmt.Printf("Receive message: %s\n", string(message))

		ToClientChan <- clientInfo{ClientId: "https://piehost.com", MessageId: "messageId", Msg: string(message)}
	}
}

func (c *Client) WriteMessage() {
	for {
		clientInfo := <-ToClientChan
		fmt.Println(clientInfo)
		if conn, err := Manager.GetByClientId(clientInfo.ClientId); err == nil && conn != nil {

			if err := conn.Socket.WriteMessage(websocket.TextMessage, []byte(clientInfo.Msg)); err != nil {
				Manager.DisConnect <- conn
				fmt.Println(err)
				return
			}
		}
	}
}

*/
