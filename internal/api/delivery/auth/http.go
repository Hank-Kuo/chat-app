package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"

	authSrv "chat-app/internal/api/service/auth"
	"chat-app/internal/dto"
	"chat-app/internal/models"
	"chat-app/pkg/logger"
	httpResponse "chat-app/pkg/response/http_response"
	"chat-app/pkg/tracer"
)

type httpHandler struct {
	authSrv authSrv.Service
	logger  logger.Logger
}

func NewHttpHandler(e *gin.RouterGroup, authSrv authSrv.Service, logger logger.Logger) {
	handler := &httpHandler{
		authSrv: authSrv,
		logger:  logger,
	}

	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)
}

func (h *httpHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "AuthHttpHandler.Register", nil)
	defer span.End()

	var body dto.RegisterReqDto
	if err := c.ShouldBindJSON(&body); err != nil {
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	if err := h.authSrv.Register(ctx, &models.User{
		Username: body.Username,
		Password: body.Password,
		Email:    body.Email,
	}); err != nil {
		tracer.AddSpanError(span, err)
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	httpResponse.OK(http.StatusOK, "register successfully", nil).ToJSON(c)
}

func (h *httpHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "AuthHttpHandler.Login", nil)
	defer span.End()

	var body dto.LoginReqDto
	if err := c.ShouldBindJSON(&body); err != nil {
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	data, err := h.authSrv.Login(ctx, body.Email, body.Password)
	if err != nil {
		tracer.AddSpanError(span, err)
		httpResponse.Fail(err, h.logger).ToJSON(c)
		return
	}

	httpResponse.OK(http.StatusOK, "login successfully", data).ToJSON(c)
}
