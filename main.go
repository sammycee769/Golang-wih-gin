package main

import (
	"log"
	"todoList/config"
	"todoList/db"
	"todoList/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.Load()
	db.Connect(cfg)
	route := gin.Default()
	routes.SetupRoutes(route)
	route.Run(":8080")

}
