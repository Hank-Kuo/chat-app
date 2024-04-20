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
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"github.com/Hank-Kuo/chat-app/config"
	"github.com/Hank-Kuo/chat-app/pkg/logger"
	"github.com/Hank-Kuo/chat-app/pkg/manager"
)

type Server struct {
	engine        *gin.Engine
	grpcEngine    *grpc.Server
	session       *gocql.Session
	snowflakeNode *snowflake.Node
	cfg           *config.Config
	db            *sqlx.DB
	rdb           *redis.Client
	manager       *manager.ClientManager
	logger        logger.Logger
}

func NewServer(cfg *config.Config, db *sqlx.DB, session *gocql.Session, rdb *redis.Client, manager *manager.ClientManager, snowflakeNode *snowflake.Node, logger logger.Logger) *Server {
	return &Server{
		engine:        nil,
		grpcEngine:    nil,
		cfg:           cfg,
		db:            db,
		rdb:           rdb,
		session:       session,
		manager:       manager,
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
