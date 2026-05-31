package orm

import (
	"fmt"
	appconfig "organization_structure/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase(config appconfig.Config) *gorm.DB {
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
