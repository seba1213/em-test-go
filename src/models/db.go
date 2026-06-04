package models

import (
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func OpenDatabaseConnection() {
	var err error
	host := os.Getenv("POSTGRES_HOST")
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	databaseName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s  TimeZone=Africa/Douala", host, username, password, databaseName, port)

	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		slog.Error("database connection failed", "error", err)
		panic(err)
	}
	slog.Info("database connected", "host", host, "dbname", databaseName, "port", port)
}

func AutoMigrateModels() {
	slog.Info("auto migrate started")

	modelsToMigrate := []struct {
		name  string
		model interface{}
	}{
		{"Subscription", &Subscription{}},
	}

	for _, item := range modelsToMigrate {
		if err := Database.AutoMigrate(item.model); err != nil {
			slog.Error("auto migrate failed", "model", item.name, "error", err)
			os.Exit(1)
		}
		slog.Info("model migrated", "model", item.name)
	}
	slog.Info("auto migrate completed")
}
