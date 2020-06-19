package config

import (
	"github.com/joho/godotenv"
	"github.com/semyon-dev/znai-krai/log"
	"os"
)

var (
	GoogleMapsAPIKey        string
	SpreadsheetIDForms      string
	SpreadsheetIDFsinPlaces string
	SpreadsheetEmail        string
	SpreadsheetCoronavirus  string
	Credentials             string
	YandexAPIKey            string
	MongoDBLogin            string
	MongoDBPass             string
	TelegramAPIToken        string
	TelegramChatID          string
)

// функция загрузки конфигов из .env файла/переменных окружения
func Load() {
	err := godotenv.Load()
	if err != nil {
		log.HandleErr(err)
	}
	GoogleMapsAPIKey = os.Getenv("GOOGLE_MAPS_API_KEY")
	SpreadsheetIDForms = os.Getenv("SPREADSHEET_ID_NEW")
	SpreadsheetIDFsinPlaces = os.Getenv("SPREADSHEET_ID_FSINPLACES")
	SpreadsheetCoronavirus = os.Getenv("SPREADSHEET_CORONAVIRUS")
	SpreadsheetEmail = os.Getenv("SPREADSHEET_EMAILS")
	Credentials = os.Getenv("CREDENTIALS_ENV")
	YandexAPIKey = os.Getenv("YANDEX_API_KEY")
	MongoDBPass = os.Getenv("MONGO_DB_PASS")
	MongoDBLogin = os.Getenv("MONGO_DB_LOGIN")
	TelegramAPIToken = os.Getenv("TELEGRAM_BOT")
	TelegramChatID = os.Getenv("CHAT_LOGS")
}
