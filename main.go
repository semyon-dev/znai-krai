package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/db"
	"github.com/semyon-dev/znai-krai/form"
	"github.com/semyon-dev/znai-krai/handlers"
	"net/http"
	"os"
)

func main() {

	// загружаем конфиги (API ключи и прочее)
	config.Load()
	gin.SetMode(os.Getenv("GIN_MODE"))

	// Подключение to Google Sheets
	handlers.Connect()

	// Подключение к MongoDB
	db.Connect()

	// обновляем места параллельно
	go handlers.UpdateAllPlaces()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"znai-krai api": "v0.9.3",
		})
	})

	// Для обработки ошибки 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "not found"})
	})

	// метод для получения всех аналитик
	router.GET("/analytics", handlers.Analytics)

	// все нарушения разом
	router.GET("/violations", handlers.Violations)

	// метод для получения всех учреждений из нашей таблицы
	router.GET("/places/:_id", handlers.Places)
	router.GET("/places/", handlers.Places)

	// метод для получения всех учреждений из нашей таблицы
	router.GET("/corona_places", handlers.CoronaPlaces)

	// отзывы с Google Maps
	router.GET("/reviews/:name", handlers.Reviews)

	// метод для создания новых форм - заявок
	router.POST("/form", handlers.NewForm)

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
