package routes

import (
	"github.com/gin-gonic/gin"
)

const apiVersion = "api/v2"

// SetupRouter sets up the router.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// users := r.Group("/users")
	// {
	// 	users.GET("/", controllers.GetUsers)
	// 	users.GET("/:id", controllers.GetUser)
	// 	users.POST("/", controllers.CreateUser)
	// 	users.PATCH("/:id", controllers.UpdateUser)
	// 	users.DELETE("/:id", controllers.DeleteUser)
	// }

	// api := r.Group(apiVersion)
	// api.Use(middlewares...) // Apply the additional middlewares passed to the function

	{
		// feed_document.RegisterHandlers(api)
		// master_data.RegisterHandlers(api)
		// activity.RegisterHandlers(api)
	}

	return r
}
