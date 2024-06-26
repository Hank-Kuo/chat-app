package http_response

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/Hank-Kuo/chat-app/pkg/customError/httpError"
	"github.com/Hank-Kuo/chat-app/pkg/logger"
)

type response struct {
	StatusCode int
	Body       *responseBody
}

type responseBody struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(statusCode int, message string, data interface{}) *response {
	return &response{
		statusCode,
		&responseBody{Status: "success", Message: message, Data: data},
	}
}

func Fail(err error, logger logger.Logger) *response {
	parseErr := httpError.ParseError(err)

	if parseErr.Detail != nil {
		logger.Error(parseErr.Detail)
	}

	return &response{
		parseErr.GetStatus(),
		&responseBody{Status: "fail", Message: parseErr.GetMessage()},
	}
}

func (r *response) ToJSON(c *gin.Context) {
	if r.Body.Status == "fail" {
		c.AbortWithStatusJSON(r.StatusCode, r.Body)
	} else {
		c.JSON(r.StatusCode, r.Body)
	}
}

func (r *response) ToWebSocketJSON(c *websocket.Conn) {
	c.WriteJSON(r.Body)
}
