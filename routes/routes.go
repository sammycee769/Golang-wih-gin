package routes

import (
	"todoList/handlers"
	"todoList/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(route *gin.Engine) {
	api := route.Group("/api/v1")
	{
		api.POST("/register", handlers.RegisterUser)
		api.POST("/login", handlers.LoginUser)
	}
		tasks := api.Group("/tasks")
		tasks.Use(middleware.RequireAuth)
		{
			tasks.GET("", handlers.GetTasks)
			tasks.POST("", handlers.CreateTasks)
			tasks.GET("/:id", handlers.GetTask)
			tasks.PUT("/:id", handlers.UpdateTask)
			tasks.DELETE("/:id", handlers.DeleteTask)
			tasks.PATCH("/:id", handlers.PatchTask)
		}
}

