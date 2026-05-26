package main

import (
	"net/http"

	"github.com/TeamPemweb/backkoshub/controllers"
	"github.com/TeamPemweb/backkoshub/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API is running",
		})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)

			authProtected := auth.Group("")
			authProtected.Use(middleware.RequireAuth)
			{
				authProtected.POST("/choose-role", controllers.ChooseRole)
				authProtected.POST("/profile/setup", controllers.SetupProfile)
				authProtected.GET("/profile", controllers.GetProfile)
			}
		}
		ownerGroup := v1.Group("/owner")
		ownerGroup.Use(middleware.RequireAuth, middleware.RoleMiddleware("pemilik"))
		{

		}

		residentGroup := v1.Group("/resident")
		residentGroup.Use(middleware.RequireAuth, middleware.RoleMiddleware("penghuni"))
		{
			
		}
	}
}