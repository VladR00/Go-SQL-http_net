package postgre

import (
	"Go-SQL-http_net/pkg/config"
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	Db *gorm.DB
}

func NewPostgreSQL(db *gorm.DB) *Storage {
	return &Storage{Db: db}
}

func ConnectPostgreSQL(cfg config.PostgreSQL) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Db, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Conns)

	slog.Info("PostgreSQL connected", "Port", cfg.Port)
	return db, nil
}
