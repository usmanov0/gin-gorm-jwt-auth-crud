package main

import (
	"log"
	"simple-crud-api/config"
	model "simple-crud-api/models"
	"simple-crud-api/storage/initializers"
)

func init() {
	config.LoadEnv()
	initializers.ConnectDb()
}

func main() {
	err := initializers.DB.Migrator().DropTable(model.User{}, model.Category{}, model.Post{}, model.Comment{})
	if err != nil {
		log.Fatal("Dropping table failed")
	}

	err = initializers.DB.AutoMigrate(model.User{}, model.Category{}, model.Post{}, model.Comment{})
	if err != nil {
		log.Fatal("migration failed")
	}
}
