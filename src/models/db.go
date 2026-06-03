package models

import (
	"fmt"
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
		panic(err)
	} else {
		fmt.Println("🚀🚀🚀🚀🚀🚀")
	}
}

func AutoMigrateModels() {
	fmt.Println("START AUTO MIGRATE MODELS")

	modelsToMigrate := []struct {
		name  string
		model interface{}
	}{
		{"Subscription", &Subscription{}},
	}

	for _, item := range modelsToMigrate {
		if err := Database.AutoMigrate(item.model); err != nil {
			fmt.Fprintf(os.Stderr, "Error in AutoMigrate for %s: %v\n", item.name, err)
			os.Exit(1)
		}
		fmt.Printf("%s model migrated successfully\n", item.name)
	}
	fmt.Println("✅ MODELS MIGRATED SUCCESSFULLY")
}
