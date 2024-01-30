package main

import (
	"github.com/gin-gonic/gin"
	"simple-crud-api/api"
	"simple-crud-api/config"
	"simple-crud-api/storage/initializers"
)

func init() {
	config.LoadEnv()
	initializers.ConnectDb()
}

func main() {
	r := gin.Default()
	api.Route(r)
	r.Run()
}
