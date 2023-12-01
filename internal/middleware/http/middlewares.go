package http

import (
	"chat-app/config"
	"chat-app/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

type Middleware struct {
	cfg    *config.Config
	logger logger.Logger
}

func NewMiddlewares(cfg *config.Config, logger logger.Logger) *Middleware {
	return &Middleware{cfg: cfg, logger: logger}
}

func NewGlobalMiddlewares(engine *gin.Engine) {
	engine.Use(corsMiddleware())
	engine.NoMethod(httpNotFound)
	engine.NoRoute(httpNotFound)
	Healthz(engine)

	m := ginmetrics.GetMonitor()

	m.SetMetricPath("/metrics")
	m.SetSlowTime(5)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	m.Use(engine)

}
