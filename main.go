package main

import (
	"os"
	"github.com/TeamPemweb/backkoshub/initializers"
	"github.com/TeamPemweb/backkoshub/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	router := gin.Default()

	router.Use(middleware.RequireAuth)
	RegisterRoutes(router)

	
	router.Run(":" + os.Getenv("PORT"))
}
