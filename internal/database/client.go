package database

import (
	"fmt"
	"log"

	"github.com/ponyo877/dummy_data_generator/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresClient() (*gorm.DB, error) {
	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("DB設定の読み込みに失敗しました: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", // sslmode=disable TimeZone=Asia/Shanghai
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Database,
		dbConfig.Port,
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if _, err := gormDB.DB(); err != nil {
		return nil, err
	}
	// defer db.Close()
	return gormDB, nil
}
