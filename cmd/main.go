package main

import (
	"github.com/gin-gonic/gin"
	"simple-crud-api/internal/common/config"
	"simple-crud-api/internal/common/db/initializers"
	"simple-crud-api/internal/ports/router"
)

func init() {
	config.LoadEnv()
	initializers.ConnectDb()
}

func main() {
	r := gin.Default()
	router.Route(r)
	r.Run()
}
