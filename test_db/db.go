package db

import (
	"github.com/joho/godotenv"
	"log"
	"simple-crud-api/models"
	"simple-crud-api/storage/initializers"
)

func DatabaseRefresh() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	initializers.ConnectDb()

	err = initializers.DB.Migrator().DropTable(models.User{}, models.Category{}, models.Post{}, models.Comment{})
	if err != nil {
		log.Fatal("Table dropping failed")
	}

	err = initializers.DB.AutoMigrate(models.User{}, models.Category{}, models.Post{}, models.Comment{})

	if err != nil {
		log.Fatal("Migration failed")
	}
}
