package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/semyon-dev/RusSeated/Sheet"
	"net/http"
)

func main() {

	gin.SetMode(gin.DebugMode)

	Sheet.Connect()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"welcome to RusSeated api": "0.1beta",
		})
	})

	// метод для создания новых форм (заявок)
	router.POST("/form", Sheet.NewForm)

	err := router.Run(":8080")
	if err != nil {
		panic(err.Error())
	}
}
