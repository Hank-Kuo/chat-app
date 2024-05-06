package server

import (
	"net/http"
	"time"

	authDelivery "github.com/Hank-Kuo/chat-app/internal/api/delivery/auth"
	channelDelivery "github.com/Hank-Kuo/chat-app/internal/api/delivery/channel"
	messageDelivery "github.com/Hank-Kuo/chat-app/internal/api/delivery/message"
	authService "github.com/Hank-Kuo/chat-app/internal/api/service/auth"
	channelService "github.com/Hank-Kuo/chat-app/internal/api/service/channel"
	messageService "github.com/Hank-Kuo/chat-app/internal/api/service/message"
	http_middleware "github.com/Hank-Kuo/chat-app/internal/middleware/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) newHttpServer(authSrv authService.Service, channelSrv channelService.Service, messageSrv messageService.Service) *http.Server {
	if s.cfg.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	http_middleware.NewGlobalMiddlewares(engine)
	middleware := http_middleware.NewMiddlewares(s.cfg, s.logger)

	api := engine.Group("/api")
	authDelivery.NewHttpHandler(api, authSrv, s.logger)
	channelDelivery.NewHttpHandler(api, channelSrv, middleware, s.logger)
	messageDelivery.NewHttpHandler(api, messageSrv, middleware, s.logger)
	messageDelivery.NewWebSocketHandler(api, s.cfg, messageSrv, channelSrv, s.manager, middleware, s.logger)

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

	go s.manager.Start()

	return httpServer
}
