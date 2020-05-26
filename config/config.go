package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var (
	GoogleMapsAPIKey        string
	SpreadsheetIDForms      string
	SpreadsheetIDFsinPlaces string
	Credentials             string
	YandexAPIKey            string
	SpreadsheetCoronavirus  string
)

// функция загрузки конфигов из .env файла/переменных окружения
func Load() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	GoogleMapsAPIKey = os.Getenv("GOOGLE_MAPS_API_KEY")
	SpreadsheetIDForms = os.Getenv("SPREADSHEET_ID")
	SpreadsheetIDFsinPlaces = os.Getenv("SPREADSHEET_ID_FSINPLACES")
	SpreadsheetCoronavirus = os.Getenv("SPREADSHEET_CORONAVIRUS")
	Credentials = os.Getenv("CREDENTIALS_ENV")
	YandexAPIKey = os.Getenv("YANDEX_API_KEY")
}
