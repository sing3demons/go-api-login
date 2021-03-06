package main

import (
	"app/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Static("/images", "./uploaded/images")

	routes.Setup(r)

	r.Run(":" + os.Getenv("PORT"))
}
