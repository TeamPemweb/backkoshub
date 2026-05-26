package main

import (
	"os"
	"github.com/TeamPemweb/backkoshub/initializers"
	"github.com/gin-gonic/gin"	
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	router := gin.Default()
	Routes(router)

	router.Run(":" + os.Getenv("PORT"))
}
