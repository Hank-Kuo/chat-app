package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"github.com/Hank-Kuo/chat-app/config"
	messageDelivery "github.com/Hank-Kuo/chat-app/internal/api/delivery/message"
	authRepository "github.com/Hank-Kuo/chat-app/internal/api/repository/auth"
	channelRepository "github.com/Hank-Kuo/chat-app/internal/api/repository/channel"
	messageRepository "github.com/Hank-Kuo/chat-app/internal/api/repository/message"
	authService "github.com/Hank-Kuo/chat-app/internal/api/service/auth"
	channelService "github.com/Hank-Kuo/chat-app/internal/api/service/channel"
	messageService "github.com/Hank-Kuo/chat-app/internal/api/service/message"
	"github.com/Hank-Kuo/chat-app/pkg/kafka"
	"github.com/Hank-Kuo/chat-app/pkg/logger"
	"github.com/Hank-Kuo/chat-app/pkg/manager"
)

type Server struct {
	engine        *gin.Engine
	grpcEngine    *grpc.Server
	session       *gocql.Session
	snowflakeNode *snowflake.Node
	kafkaProducer *kafka.KafkaProducer
	cfg           *config.Config
	db            *sqlx.DB
	rdb           *redis.Client
	manager       *manager.ClientManager
	logger        logger.Logger
}

func NewServer(cfg *config.Config, db *sqlx.DB, session *gocql.Session, rdb *redis.Client, manager *manager.ClientManager, snowflakeNode *snowflake.Node, kafkaProducer *kafka.KafkaProducer, logger logger.Logger) *Server {
	return &Server{
		engine:        nil,
		grpcEngine:    nil,
		cfg:           cfg,
		db:            db,
		rdb:           rdb,
		kafkaProducer: kafkaProducer,
		session:       session,
		manager:       manager,
		snowflakeNode: snowflakeNode,
		logger:        logger,
	}
}

func (s *Server) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	authRepo := authRepository.NewRepo(s.db)
	authSrv := authService.NewService(s.cfg, authRepo, s.logger)
	channelRepo := channelRepository.NewRepo(s.db, s.rdb)
	channelSrv := channelService.NewService(s.cfg, channelRepo, s.logger)
	messageRepo := messageRepository.NewRepo(s.session, s.kafkaProducer, s.rdb)
	messageSrv := messageService.NewService(s.cfg, messageRepo, s.snowflakeNode, s.logger)

	// http handler
	httpServer := s.newHttpServer(authSrv, channelSrv, messageSrv)

	// kafka consumer handler
	adminClient, err := kafka.NewKafkaAdmin(s.cfg.Kafka)
	defer adminClient.Close()
	client, err := kafka.NewKafkaConsumer(s.cfg.Kafka)
	defer client.Consumer.Close()

	messageKafkaHandler := messageDelivery.NewKafkaHandler(s.cfg, messageSrv, channelSrv, s.manager, s.logger)
	go messageKafkaHandler.Listen(client, adminClient)
	go messageKafkaHandler.ConsumeMessage(client)

	// grpc handler
	grpcClose, grpcServer, err := s.newGrpcServer()
	if err != nil {
		s.logger.Fatalf("Error gprc serve: %s", err)
	}
	defer grpcClose()

	// graceful shutdown
	<-ctx.Done()
	s.logger.Info("Shutdown Server ...")

	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(ctx); err != nil {
		s.logger.Fatal(err)
	}

}
