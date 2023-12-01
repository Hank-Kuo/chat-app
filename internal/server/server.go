package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	// "github.com/segmentio/kafka-go"
	"github.com/gocql/gocql"
	"google.golang.org/grpc"

	"chat-app/config"
	"chat-app/pkg/logger"
)

type Server struct {
	engine        *gin.Engine
	grpcEngine    *grpc.Server
	session       *gocql.Session
	snowflakeNode *snowflake.Node
	cfg           *config.Config
	db            *sqlx.DB
	logger        logger.Logger
}

func NewServer(cfg *config.Config, db *sqlx.DB, session *gocql.Session, snowflakeNode *snowflake.Node, logger logger.Logger) *Server {
	return &Server{
		engine: nil, grpcEngine: nil,
		cfg: cfg, db: db,
		session:       session,
		snowflakeNode: snowflakeNode,
		logger:        logger,
	}
}

func (s *Server) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	httpServer := s.newHttpServer()

	// graceful shutdown
	<-ctx.Done()
	s.logger.Info("Shutdown Server ...")

	if err := httpServer.Shutdown(ctx); err != nil {
		s.logger.Fatal(err)
	}
}
