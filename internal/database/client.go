package database

import (
	"fmt"
	"log"

	"github.com/ponyo877/dummy_data_generator/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabaseClient(dbEngine string) (*gorm.DB, error) {
	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("failed to load DB config: %v", err)
	}
	var gormDB *gorm.DB
	switch dbEngine {
	case "postgres":
		gormDB, err = PostgresClient(dbConfig)
	case "mysql":
		gormDB, err = MySQLClient(dbConfig)
	default:
		log.Fatalf("%s is not supported", dbEngine)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	// db, err := gormDB.DB()
	// if err != nil {
	// 	return nil, err
	// }
	// defer db.Close()
	return gormDB, nil
}

func PostgresClient(config config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", // sslmode=disable TimeZone=Asia/Shanghai
		config.Host,
		config.User,
		config.Password,
		config.Database,
		config.Port,
	)
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}

func MySQLClient(config config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
