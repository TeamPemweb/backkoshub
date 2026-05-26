package main

import (
	"net/http"

	"github.com/TeamPemweb/backkoshub/controllers"
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
    }
	}
}