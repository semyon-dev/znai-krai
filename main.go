package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/semyon-dev/znai-krai/Sheet"
	"net/http"
)

func main() {

	gin.SetMode(gin.DebugMode)

	Sheet.Connect()

	go Sheet.UpdatePlaces()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Welcome to RusSeated api": "",
		})
	})

	// метод для получения всех учреждений из нашей таблицы
	router.GET("/places", Sheet.Places)

	// отзывы Google Maps
	router.GET("/reviews/:name", Sheet.Reviews)

	// метод для создания новых форм (заявок)
	router.POST("/form", Sheet.NewForm)

	err := router.Run(":8080")
	if err != nil {
		panic(err.Error())
	}
}
