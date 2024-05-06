package server

import (
	"net"
	"time"

	messageDelivery "github.com/Hank-Kuo/chat-app/internal/api/delivery/message"
	grpc_middleware "github.com/Hank-Kuo/chat-app/internal/middleware/grpc"
	messagePb "github.com/Hank-Kuo/chat-app/pb/message"

	go_grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	go_grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	go_grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	go_grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	go_grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

func (s *Server) newGrpcServer() (func() error, *grpc.Server, error) {
	l, err := net.Listen("tcp", ":"+s.cfg.Server.GrpcPort)
	if err != nil {
		return nil, nil, errors.Wrap(err, "grpc.net.Listen")
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		grpc.UnaryInterceptor(go_grpc_middleware.ChainUnaryServer(
			go_grpc_ctxtags.UnaryServerInterceptor(),
			go_grpc_opentracing.UnaryServerInterceptor(),
			go_grpc_prometheus.UnaryServerInterceptor,
			go_grpc_recovery.UnaryServerInterceptor(),
			grpc_middleware.LoggingInterceptor(s.logger),
		)),
	)

	go_grpc_prometheus.Register(grpcServer)

	messageHandler := messageDelivery.NewGrpcHandler(s.manager, s.logger)
	messagePb.RegisterMessageServiceServer(grpcServer, messageHandler)

	if s.cfg.Server.Debug {
		reflection.Register(grpcServer)
	}

	go func() {
		if err := grpcServer.Serve(l); err != nil {
			s.logger.Fatalf("Error gprc serve: %s", err)
		}
	}()

	return l.Close, grpcServer, nil
}
