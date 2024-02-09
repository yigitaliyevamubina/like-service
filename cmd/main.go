package main

import (
	"google.golang.org/grpc"
	"like-service/config"
	pb "like-service/genproto/like_service"
	"like-service/pkg/db"
	"like-service/pkg/logger"
	"like-service/service"
	grpcClient "like-service/service/grpc_client"
	"net"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "like-service")
	defer logger.Cleanup(log)

	log.Info("main: sqlConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostServicePort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, _, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sql connection to postgres error", logger.Error(err))
	}

	client, err := grpcClient.New(cfg)
	if err != nil {
		log.Fatal("error while adding grpc client", logger.Error(err))
	}

	likeService := service.NewLikeService(connDB, log, client)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("failed to listen to: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterLikeServiceServer(s, likeService)
	log.Info("main: server is runnning",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("error while listening: %v", logger.Error(err))
	}
}
