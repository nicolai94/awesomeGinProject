package router

import (
	"awesomeProject/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(init *config.Initialization) *gin.Engine {

	router := gin.New()
	router.Use(cors.Default())

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")

	api.GET("/users", init.UserCtrl.GetAllUserData)
	{
		user := api.Group("/user")
		user.POST("", init.UserCtrl.AddUserData)
		user.GET("/:userID", init.UserCtrl.GetUserById)
		user.PUT("/:userID", init.UserCtrl.UpdateUserData)
		user.DELETE("/:userID", init.UserCtrl.DeleteUser)
	}

	{
		auth := api.Group("/auth")
		auth.POST("/login", init.AuthCtrl.Login)
	}
	return router
}
