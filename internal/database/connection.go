package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ruslannnnnnnnn/test-file-storage/internal/common"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabaseConnection(dbConfig common.DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s", dbConfig.Host, dbConfig.User, dbConfig.DbName, strconv.Itoa(dbConfig.Port))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db

}
