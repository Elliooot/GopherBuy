package repository

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"GopherBuy/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var once sync.Once

func InitDB() error {
	var err error
	once.Do(func() {
		// dsn := "host=localhost user=gorm password=gorm dbname=gopherBuy port=5432 sslmode=disable TimeZone=Asia/Taipei"
		host := os.Getenv("POSTGRES_HOST")
		if host == "" {
			host = "localhost"
		}

		port := os.Getenv("POSTGRES_PORT")
		user := os.Getenv("POSTGRES_USER")
		pw := os.Getenv("POSTGRES_PASSWORD")
		dbname := os.Getenv("POSTGRES_DB")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, pw, dbname, port)

		// Dealing with the status that gorm connecting to sqlDB
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Cannot connect to database: %v", err)
		}

		// type DB is a defined struct => https://pkg.go.dev/database/sql#DB
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal(err)
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		DB = db

		db.AutoMigrate(&model.User{}, &model.Product{}, &model.Order{}, &model.FlashSale{})

		fmt.Println("Database has been initialised and migrated")
	})
	return err
}
