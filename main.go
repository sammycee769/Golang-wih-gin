package main

import (
	"todoList/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin.Default() gives us a router with Logger and Recovery middleware
	route := gin.Default()
	routes.SetupRoutes(route)
	route.Run(":8080")

}
