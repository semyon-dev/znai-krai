package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/db"
	"github.com/semyon-dev/znai-krai/form"
	"github.com/semyon-dev/znai-krai/handlers"
	log2 "github.com/semyon-dev/znai-krai/log"
	"github.com/semyon-dev/znai-krai/sheet"
	"net/http"
	"os"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// загружаем конфиги (API ключи и прочее)
	config.Load()
	log2.ConnectBot()

	gin.SetMode(os.Getenv("GIN_MODE"))

	// Подключение to Google Sheets
	sheet.Connect()

	// Подключение к MongoDB
	db.Connect()

	// обновляем места параллельно
	go handlers.UpdateAllPlaces()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"znai-krai api": "v0.14.0",
		})
	})

	// Для обработки ошибки 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "not found"})
	})

	// метод для получения аналитики
	router.GET("/analytics", handlers.Analytics)

	// пояснение - перевод
	router.GET("/explanations", handlers.Explanations)

	// все нарушения разом
	router.GET("/violations", handlers.Violations)

	// методы для получения учреждений из нашей таблицы
	router.GET("/places/:_id", handlers.Places)
	router.GET("/places/", handlers.Places)

	// метод для получения всех учреждений из нашей таблицы
	router.GET("/corona_places", handlers.CoronaPlaces)

	// Deprecated: отзывы с Google Maps
	router.GET("/reviews/:name", handlers.Reviews)

	// метод для создания новых форм - заявок
	router.POST("/form", handlers.NewForm)
	router.POST("/form_corona", handlers.NewFormCorona)

	// получение всех вопросов для заполнения со стороны клиента
	router.GET("/formQuestions", form.Questions)

	// репорт для ошибок/багов
	router.POST("/report", form.Report)

	// подписка на почтовую рассылку
	router.POST("/mailing", handlers.NewMailing)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	err := router.Run(":" + port)
	if err != nil {
		log.Panic().AnErr("Не получилось запустить:", err)
	}
}
