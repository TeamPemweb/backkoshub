package main

import (
	"os"
	"github.com/TeamPemweb/backkoshub/initializers"
	"github.com/gin-gonic/gin"	
	"github.com/TeamPemweb/backkoshub/middleware"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	Routes(router)

	router.Run(":" + os.Getenv("PORT"))
}
