package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"organization_structure/internal/adapters/database"
	"organization_structure/internal/adapters/protocol"
	appconfig "organization_structure/internal/config"
	"organization_structure/internal/logging"
	"organization_structure/internal/usecase"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := appconfig.NewConfig()

	logging.InitLogger(&config)

	db := setupDatabase(config)
	p := database.NewPgStorage(db)
	u := usecase.NewUseCase(p)
	r := protocol.NewRouter(u)

	var server = &http.Server{
		Addr:    ":" + strconv.Itoa(config.Server.Port),
		Handler: r,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	<-mainCtx.Done()

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}

	log.Println("shutting down database...")

	sqlDB, _ := db.DB()

	if err := sqlDB.Close(); err != nil {
		log.Fatal(err)
	}
}

func setupDatabase(config appconfig.Config) *gorm.DB {
	var attempts = 5
	var delay = 2 * time.Second

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Db.Host, config.Db.User, config.Db.Password, config.Db.Name, config.Db.Port, config.Db.SslMode)

	for i := 0; i < attempts; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: true})

		if err == nil {
			sqlDB, _ := db.DB()

			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Minute * 5)

			return db
		}

		time.Sleep(delay)
	}

	panic("failed to connect to database")
}
