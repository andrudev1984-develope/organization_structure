package main

import (
	"context"
	"log"
	"net/http"
	"organization_structure/internal/adapters/database"
	"organization_structure/internal/adapters/orm"
	"organization_structure/internal/adapters/protocol"
	appconfig "organization_structure/internal/config"
	"organization_structure/internal/logging"
	"organization_structure/internal/usecase"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := appconfig.NewConfig()

	logging.InitLogger(&config)

	db := orm.SetupDatabase(config)
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
