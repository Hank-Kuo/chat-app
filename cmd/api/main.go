package main

import (
	"context"
	"fmt"
	"log"

	"chat-app/config"
	"chat-app/internal/server"
	"chat-app/pkg/database"
	"chat-app/pkg/logger"
	"chat-app/pkg/tracer"

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

	// kakfaWriter, err := kafka.NewWriter(cfg.Kafka, "user_email")
	// if err != nil {
	// 	apiLogger.Fatal("Can't connect with kafka", err)
	// }
	// defer kakfaWriter.Close()

	cassandraSess, err := database.ConnectCassandra(&cfg.Cassandra)

	if err != nil {
		panic(fmt.Errorf("load cassandra: %v", err))
	}

	defer cassandraSess.Close()

	snowflakeNode, err := snowflake.NewNode(1)
	if err != nil {
		panic(fmt.Errorf("load snowflake: %v", err))
	}

	// init server
	srv := server.NewServer(cfg, db, cassandraSess, snowflakeNode, apiLogger)
	srv.Run()
}
