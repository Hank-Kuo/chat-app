package server

import (
	"net/http"
	"time"

	authDelivery "chat-app/internal/api/delivery/auth"
	channelDelivery "chat-app/internal/api/delivery/channel"
	messageDelivery "chat-app/internal/api/delivery/message"
	authRepository "chat-app/internal/api/repository/auth"
	channelRepository "chat-app/internal/api/repository/channel"
	messageRepository "chat-app/internal/api/repository/message"
	authService "chat-app/internal/api/service/auth"
	channelService "chat-app/internal/api/service/channel"
	messageService "chat-app/internal/api/service/message"
	http_middleware "chat-app/internal/middleware/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerHttpHanders(engine *gin.Engine) {
	middleware := http_middleware.NewMiddlewares(s.cfg, s.logger)

	api := engine.Group("/api")

	authRepo := authRepository.NewRepo(s.db)
	authSrv := authService.NewService(s.cfg, authRepo, s.logger)
	authDelivery.NewHttpHandler(api, authSrv, s.logger)

	channelRepo := channelRepository.NewRepo(s.db)
	channelSrv := channelService.NewService(s.cfg, channelRepo, s.logger)
	channelDelivery.NewHttpHandler(api, channelSrv, middleware, s.logger)

	messageRepo := messageRepository.NewRepo(s.session)
	messageSrv := messageService.NewService(s.cfg, messageRepo, s.snowflakeNode, s.logger)
	messageDelivery.NewHttpHandler(api, messageSrv, s.logger)

}

func (s *Server) newHttpServer() *http.Server {
	if s.cfg.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	http_middleware.NewGlobalMiddlewares(engine)

	s.registerHttpHanders(engine)

	httpServer := &http.Server{
		Addr:           ":" + s.cfg.Server.Port,
		Handler:        engine,
		ReadTimeout:    time.Second * time.Duration(s.cfg.Server.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(s.cfg.Server.WriteTimeout),
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Error http ListenAndServe: %s", err)
		}
	}()

	return httpServer
}
