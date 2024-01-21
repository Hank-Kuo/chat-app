package message

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	messageSrv "chat-app/internal/api/service/message"
	"chat-app/internal/dto"
	"chat-app/pkg/logger"
	httpResponse "chat-app/pkg/response/http_response"
	"chat-app/pkg/tracer"
)

type httpHandler struct {
	messageSrv messageSrv.Service
	logger     logger.Logger
}

func NewHttpHandler(e *gin.RouterGroup, messageSrv messageSrv.Service, logger logger.Logger) {
	handler := &httpHandler{
		messageSrv: messageSrv,
		logger:     logger,
	}
	e.GET("/message", handler.GetMessage) // get message
	e.GET("/reply", handler.GetReply)     // get reply

	e.POST("/message", handler.CreateMessage) // send message
	e.POST("/reply", handler.CreateReply)     // send reply

}

func (h *httpHandler) CreateMessage(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "MessageHttpHandler.CreateMessage", nil)
	defer span.End()

	var body dto.CreateMessageReqDto
	if err := c.ShouldBindJSON(&body); err != nil {
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	message, err := h.messageSrv.CreateMessage(ctx, &body)
	if err != nil {
		tracer.AddSpanError(span, err)
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	httpResponse.OK(http.StatusOK, "create message successfully", message).ToJSON(c)

}

func (h *httpHandler) CreateReply(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "MessageHttpHandler.CreateReply", nil)
	defer span.End()

	var body dto.CreateReplyReqDto
	if err := c.ShouldBindJSON(&body); err != nil {
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	reply, err := h.messageSrv.CreateReply(ctx, &body)
	if err != nil {
		tracer.AddSpanError(span, err)
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	httpResponse.OK(http.StatusOK, "create reply successfully", reply).ToJSON(c)
}

func (h *httpHandler) GetMessage(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "MessageHttpHandler.GetMessage", nil)
	defer span.End()

	var queryParams dto.GetMessageQueryDto
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	messages, err := h.messageSrv.GetMessage(ctx, &queryParams)
	if err != nil {
		tracer.AddSpanError(span, err)
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	fmt.Println(messages)
	httpResponse.OK(http.StatusOK, "get message successfully", messages).ToJSON(c)
}

func (h *httpHandler) GetReply(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "MessageHttpHandler.GetReply", nil)
	defer span.End()

	var queryParams dto.GetReplyQueryDto
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	reply, err := h.messageSrv.GetReply(ctx, &queryParams)
	if err != nil {
		tracer.AddSpanError(span, err)

		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	httpResponse.OK(http.StatusOK, "get reply successfully", reply).ToJSON(c)
}
