package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var (
	GoogleMapsAPIKey        string
	SpreadsheetID           string
	SpreadsheetIDFsinPlaces string
	Credentials             string
	YandexAPIKey            string
)

// функция загрузки конфигов из .env файла/переменных окружения
func Load() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	GoogleMapsAPIKey = os.Getenv("GOOGLE_MAPS_API_KEY")
	SpreadsheetID = os.Getenv("SPREADSHEET_ID")
	SpreadsheetIDFsinPlaces = os.Getenv("SPREADSHEET_ID_FSINPLACES")
	Credentials = os.Getenv("CREDENTIALS_ENV")

	YandexAPIKey = os.Getenv("YANDEX_API_KEY")
}
