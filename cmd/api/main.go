package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Hank-Kuo/chat-app/config"
	"github.com/Hank-Kuo/chat-app/internal/server"
	"github.com/Hank-Kuo/chat-app/pkg/database"
	"github.com/Hank-Kuo/chat-app/pkg/kafka"
	"github.com/Hank-Kuo/chat-app/pkg/logger"
	"github.com/Hank-Kuo/chat-app/pkg/manager"
	"github.com/Hank-Kuo/chat-app/pkg/tracer"

	"github.com/bwmarrin/snowflake"
)

func main() {
	log.Println("Starting chat-app server")
	cfg, err := config.GetConf()

	if err != nil {
		panic(fmt.Errorf("load config: %v", err))
	}

	apiLogger := logger.NewApiLogger(cfg)
	apiLogger.InitLogger()

	// init database
	db, err := database.ConnectDB(&cfg.Database)
	if err != nil {
		panic(fmt.Errorf("load database: %v", err))
	}
	defer db.Close()

	traceProvider, err := tracer.NewJaeger(cfg)
	if err != nil {
		apiLogger.Fatal("Cannot create tracer", err)
	} else {
		apiLogger.Info("Jaeger connected")
	}

	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			apiLogger.Error("Cannot shutdown tracer", err)
		}
	}()

	cassandraSess, err := database.ConnectCassandra(&cfg.Cassandra)

	if err != nil {
		panic(fmt.Errorf("load cassandra: %v", err))
	}

	defer cassandraSess.Close()

	snowflakeNode, err := snowflake.NewNode(cfg.Server.InstanceID)
	if err != nil {
		panic(fmt.Errorf("load snowflake: %v", err))
	}
	kafkaMessageWriter, err := kafka.NewWriter(cfg.Kafka, "message_topic")
	if err != nil {
		panic(fmt.Errorf("Can't connect with messsage kafka: %v", err))
	}
	defer kafkaMessageWriter.Close()

	kafkaReplyWriter, err := kafka.NewWriter(cfg.Kafka, "reply_topic")
	if err != nil {
		panic(fmt.Errorf("Can't connect with reply kafka: %v", err))

	}
	defer kafkaReplyWriter.Close()

	rdb := database.ConnectRedis(&cfg.Redis)

	manager := manager.NewClientManager(rdb, cfg.Server.InstanceIP)

	// init server
	srv := server.NewServer(cfg, db, cassandraSess, rdb, manager, snowflakeNode, kafkaMessageWriter, kafkaReplyWriter, apiLogger)
	srv.Run()
}
