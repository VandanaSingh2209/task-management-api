package Routes

import (
	"golang-assesment/Controllers"
	"golang-assesment/Middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	public := r.Group("/public")
	{
		public.GET("/tasklist", Controllers.GetTaskList)
		public.POST("/login", Controllers.LoginAuth)
	}

	protected := r.Group("/protected")
	protected.Use(Middleware.ValidateToken())
	{
		protected.Use(Middleware.RateLimitMiddleware(5, 10))
		protected.GET("/tasks", Controllers.GetTask)
		protected.POST("/create-task", Controllers.CreateTask)
		protected.PUT("/tasks/:id", Controllers.UpdateTask)
		protected.DELETE("/tasks/:id", Controllers.DeleteTask)
	}
}
