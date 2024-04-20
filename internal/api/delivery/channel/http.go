package channel

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	channelSrv "github.com/Hank-Kuo/chat-app/internal/api/service/channel"
	"github.com/Hank-Kuo/chat-app/internal/dto"
	httpMiddleware "github.com/Hank-Kuo/chat-app/internal/middleware/http"
	"github.com/Hank-Kuo/chat-app/internal/models"
	"github.com/Hank-Kuo/chat-app/pkg/logger"
	httpResponse "github.com/Hank-Kuo/chat-app/pkg/response/http_response"
	"github.com/Hank-Kuo/chat-app/pkg/tracer"
)

type httpHandler struct {
	channelSrv channelSrv.Service
	logger     logger.Logger
}

func NewHttpHandler(e *gin.RouterGroup, channelSrv channelSrv.Service, mid *httpMiddleware.Middleware, logger logger.Logger) {
	handler := &httpHandler{
		channelSrv: channelSrv,
		logger:     logger,
	}

	e.GET("/channel", mid.AuthMiddleware(), handler.Get)
	e.GET("/user/channel", mid.AuthMiddleware(), handler.GetUserChannel)
	e.POST("/channel", mid.AuthMiddleware(), handler.Create)
	e.POST("/channel/join", mid.AuthMiddleware(), handler.Join)
}

func (h *httpHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "ChannelHttpHandler.Get", nil)
	defer span.End()

	channels, err := h.channelSrv.Get(ctx)
	if err != nil {
		tracer.AddSpanError(span, err)
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}
	httpResponse.OK(http.StatusOK, "get channel successfully", channels).ToJSON(c)
}

func (h *httpHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "ChannelHttpHandler.Create", nil)
	defer span.End()

	var body dto.CreateChannelReqDto
	if err := c.ShouldBindJSON(&body); err != nil {
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	if err := h.channelSrv.Create(ctx, &models.Channel{
		Name:   body.Name,
		UserID: fmt.Sprintf("%v", userID),
	}); err != nil {
		tracer.AddSpanError(span, err)
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	httpResponse.OK(http.StatusOK, "create channel successfully", nil).ToJSON(c)
}

func (h *httpHandler) Join(c *gin.Context) {
	userID, _ := c.Get("userID")

	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "ChannelHttpHandler.Join", nil)
	defer span.End()

	var body dto.JoinChannelReqDto
	if err := c.ShouldBindJSON(&body); err != nil {
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	if err := h.channelSrv.Join(ctx, fmt.Sprintf("%v", userID), body.ChannelID); err != nil {
		tracer.AddSpanError(span, err)
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	httpResponse.OK(http.StatusOK, "user join channel successfully", nil).ToJSON(c)
}

func (h *httpHandler) GetUserChannel(c *gin.Context) {

	userID, _ := c.Get("userID")

	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "ChannelHttpHandler.GetUserChannel", nil)
	defer span.End()

	channels, err := h.channelSrv.GetUserChannel(ctx, fmt.Sprintf("%v", userID))
	if err != nil {
		tracer.AddSpanError(span, err)
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}
	httpResponse.OK(http.StatusOK, "get user channels successfully", channels).ToJSON(c)
}
