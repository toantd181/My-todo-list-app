package routes

import (
	"database/sql"
	"todo-app/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	todoHandler := handlers.NewTodoHandler(db)

	api := router.Group("/api")
	{
		todos := api.Group("/todos")
		{
			todos.GET("", todoHandler.GetTodos)
			todos.POST("", todoHandler.CreateTodo)
			todos.PUT("/:id", todoHandler.UpdateTodo)
			todos.DELETE("/:id", todoHandler.DeleteTodo)
			todos.POST("/reorder", todoHandler.ReorderTodos)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Todo App API is running",
		})
	})
}
