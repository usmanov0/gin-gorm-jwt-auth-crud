package main

import (
	"simple-crud-api/api"
	"simple-crud-api/config"
	"simple-crud-api/storage/initializers"
)

func init() {
	config.LoadEnv()
	initializers.ConnectDb()
}

func main() {
	r := api.Route()
	r.Run()
}
