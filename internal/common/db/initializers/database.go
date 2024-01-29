package initializers

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectDb() {
	connectionString := os.Getenv("DNS")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")

	defer db.Close()
}
