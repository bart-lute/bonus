package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// trusted proxies
	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()

	router.GET("/ping", ping)
	router.POST("/list", list)
	_ = router.Run(fmt.Sprintf(":%s", os.Getenv("GIN_PORT")))
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func list(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "foobar",
	})
}
