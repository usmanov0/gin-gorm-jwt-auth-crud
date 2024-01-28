package main

import (
	"log"
	"simple-crud-api/internal/common/config"
	"simple-crud-api/internal/common/db/initializers"
	"simple-crud-api/internal/models"
)

func init() {
	config.LoadEnv()
	initializers.ConnectDb()
}

func main() {
	err := initializers.DB.Migrator().DropTable(models.User{}, models.Post{}, models.Category{}, models.Comment{})
	if err != nil {
		log.Fatal("Dropping table failed")
	}

	err = initializers.DB.AutoMigrate(models.User{}, models.Post{}, models.Category{}, models.Comment{})
	if err != nil {
		log.Fatal("migration failed")
	}
}
