package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/form"
	"github.com/semyon-dev/znai-krai/sheet"
	"net/http"
	"os"
)

func main() {

	// загружаем конфиги (API ключи и прочее)
	config.Load()

	gin.SetMode(gin.DebugMode)

	// Connect to Google Sheets
	sheet.Connect()

	// обновляем места параллельно
	go sheet.UpdatePlaces()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"znai-krai api": "v0.3",
		})
	})

	// метод для получения всех учреждений из нашей таблицы
	router.GET("/places", sheet.Places)

	// отзывы с Google Maps
	router.GET("/reviews/:name", sheet.Reviews)

	// метод для создания новых форм - заявок
	router.POST("/form", sheet.NewForm)

	// получение всех вопросов для заполнения со стороны клиента
	router.GET("/formQuestions", form.Questions)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	err := router.Run(":" + port)
	if err != nil {
		panic("Не получилось запустить:" + err.Error())
	}
}
