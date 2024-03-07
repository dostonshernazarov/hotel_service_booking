package main

import (
	"fmt"
	"google.golang.org/grpc"
	"hotel_service_booking/config"
	pb "hotel_service_booking/genproto/hotel_proto"
	"hotel_service_booking/pkg/db"
	"hotel_service_booking/pkg/logger"
	"hotel_service_booking/queue/rabbitmq/consumermq"
	"hotel_service_booking/service"
	"net"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "template-service")
	defer logger.Cleanup(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	// rabbit mq -------------
	consumer, err := consumermq.NewRabbitMQConsumer("amqp://guest:guest@localhost:5672/", "test-topic")
	if err != nil {
		log.Error("NewRabbitMqConsumer", logger.Error(err))
		return
	}
	defer consumer.Close()

	go func() {
		consumer.ConsumerMassages(consumerHandler)
	}()
	// rabbit mq end ------------

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	userService := service.NewUserService(connDB, log)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterHotelServiceServer(s, userService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}

func consumerHandler(message []byte) {
	fmt.Println("Consumer 2:", string(message))
}
