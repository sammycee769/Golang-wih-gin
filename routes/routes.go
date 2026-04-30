package routes

import (
	"todoList/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(route *gin.Engine) {
	api := route.Group("/api/v1")
	{
		tasks := api.Group("/tasks")
		{
			tasks.GET("", handlers.GetTasks)
			tasks.POST("", handlers.CreateTasks)
			tasks.GET("/:id", handlers.GetTask)
			tasks.PUT("/:id", handlers.UpdateTask)
			tasks.DELETE("/:id", handlers.DeleteTask)
		}
	}
}
