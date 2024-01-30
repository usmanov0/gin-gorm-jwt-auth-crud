package main

import (
	"log"
	"simple-crud-api/config"
	models2 "simple-crud-api/models"
	"simple-crud-api/storage/initializers"
)

func init() {
	config.LoadEnv()
	initializers.ConnectDb()
}

func main() {
	err := initializers.DB.Migrator().DropTable(models2.User{}, models2.Post{}, models2.Category{}, models2.Comment{})
	if err != nil {
		log.Fatal("Dropping table failed")
	}

	err = initializers.DB.AutoMigrate(models2.User{}, models2.Post{}, models2.Category{}, models2.Comment{})
	if err != nil {
		log.Fatal("migration failed")
	}
}
