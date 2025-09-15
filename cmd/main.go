package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	pb "github.com/ruslannnnnnnnn/test-file-storage/api/gen/go/service/v1"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/api/grpc_controller"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/common"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/database"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/repository"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/service"
	"google.golang.org/grpc"
)

const envPath = "/app/.env"

func main() {

	// переменные окружения
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// подключение к бд
	dbPort := os.Getenv("POSTGRES_PORT")
	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		log.Fatal("Error parsing POSTGRES_PORT")
	}

	dbConfig := common.DatabaseConfig{
		Host:   os.Getenv("POSTGRES_HOST"),
		Port:   dbPortInt,
		User:   os.Getenv("POSTGRES_USER"),
		DbName: os.Getenv("POSTGRES_DB"),
	}

	dbConn := database.GetDatabaseConnection(dbConfig)

	// миграция таблиц
	fileRepo := repository.NewFileRepository(dbConn)
	err = fileRepo.AutoMigrate()
	if err != nil {
		log.Fatal(err)
	}

	// сервис

	fileService := service.NewFileService(fileRepo)

	// grpc сервер и контроллер
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// контроллер
	fileController := grpc_controller.NewGrpcFileController(fileService)
	pb.RegisterFileServiceServer(grpcServer, fileController)
	//

	log.Printf("gRPC server starting on port 8080")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
