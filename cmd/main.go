package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	pb "github.com/ruslannnnnnnnn/test-file-storage/api/gen/go/service/v1"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/api/grpc_controller"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/common"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/database"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/limiter"
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

	// задаём лимиты
	listFilesMaxConcurrentRequests, err := strconv.Atoi(os.Getenv("LIST_FILES_MAX_CONCURRENT_REQUESTS"))
	if err != nil {
		log.Fatal(err)
	}
	readwriteFilesMaxConcurrentRequests, err := strconv.Atoi(os.Getenv("READWRITE_FILES_MAX_CONCURRENT_REQUESTS"))
	if err != nil {
		log.Fatal(err)
	}

	limits := map[string]int{
		"/service.v1.FileService/Upload":    readwriteFilesMaxConcurrentRequests,
		"/service.v1.FileService/Download":  readwriteFilesMaxConcurrentRequests,
		"/service.v1.FileService/ListFiles": listFilesMaxConcurrentRequests,
	}

	// создаём лимитер с TTL (например, 5 минут для очистки клиентов)
	reqLimiter := limiter.NewLimiter(limits, 5*time.Minute)
	defer reqLimiter.Stop()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(reqLimiter.UnaryInterceptor),
		grpc.StreamInterceptor(reqLimiter.StreamInterceptor),
	)

	// контроллер
	fileController := grpc_controller.NewGrpcFileController(fileService)
	pb.RegisterFileServiceServer(grpcServer, fileController)
	//

	log.Printf("gRPC server starting on port 8080")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
