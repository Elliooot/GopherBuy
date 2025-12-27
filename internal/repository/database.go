package repository

import (
	"fmt"
	"log"
	"time"

	"GopherBuy/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Taipei"
	// Dealing with the status that gorm connecting to sqlDB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to database: %v", err)
	}

	// type DB is a defined struct => https://pkg.go.dev/database/sql#DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db.AutoMigrate(&model.User{}, &model.Product{}, &model.Order{})

	fmt.Printf("Database has been initialised and migrated")
}
