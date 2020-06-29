package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/db"
	"github.com/semyon-dev/znai-krai/form"
	"github.com/semyon-dev/znai-krai/handlers"
	mylog "github.com/semyon-dev/znai-krai/log"
	"github.com/semyon-dev/znai-krai/sheet"
	"net/http"
	"time"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// загружаем конфиги (API ключи и прочее)
	config.Load()
	mylog.ConnectBot()

	gin.SetMode(config.GinMode)

	// Подключение to Google Sheets
	sheet.Connect()

	// Подключение к MongoDB
	db.Connect()

	// обновляем места параллельно
	go handlers.UpdateAllPlaces()

	router := gin.Default()
	router.Use(cors.Default())
	private := cors.Default()
	if config.GinMode == config.GinModeRelease {
		private = cors.New(cors.Config{
			AllowOrigins:
			[]string{
				"https://znaikrai.herokuapp.com",
				"http://znaikrai.herokuapp.com",
				"https://znai-krai.zekovnet.ru",
				"https://znay-kray.zekovnet.ru",
				"https://znai-krai.zekovnet.ru",
				"https://znaj-kraj.zekovnet.ru",
				"http://znaj-kraj.zekovnet.ru",
				"http://znai-krai.zekovnet.ru",
				"http://znay-kray.zekovnet.ru",
				"http://znai-krai.zekovnet.ru",
			},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Host"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"znai-krai API": config.APIVersion,
		})
	})

	// robots.txt file
	router.StaticFile("/robots.txt", "./robots.txt")

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

	// получение всех вопросов для заполнения со стороны клиента
	router.GET("/formQuestions", form.Questions)

	router.Use(private)

	// Deprecated: отзывы с Google Maps
	router.GET("/reviews/:name", handlers.Reviews).Use(private)

	// метод для создания новых форм - заявок
	router.POST("/form", handlers.NewForm).Use(private)
	router.POST("/form_corona", handlers.NewFormCorona).Use(private)

	// репорт для ошибок/багов
	router.POST("/report", form.Report).Use(private)

	// подписка на почтовую рассылку
	router.POST("/mailing", handlers.NewMailing).Use(private)

	err := router.Run(":" + config.Port)
	if err != nil {
		mylog.HandlePanicWitMsg("Не получилось запустить: ", err)
	}
}
