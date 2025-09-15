package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/common"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/database"
	"google.golang.org/grpc"
)

const envPath = "/app/.env"

func main() {

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	_ = database.GetDatabaseConnection(dbConfig)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	log.Printf("gRPC server starting on port 8080")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
