// Здесь содержатся вспомогательные
// функции для анализа и сбора данных которые
// я написал для временного или единоразово использования
// В данным момент эти функции не используются и не поддерживаются, но могут понадобиться.
// Deprecated:
package util

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/model"
	"github.com/tidwall/gjson"
	"googlemaps.github.io/maps"
	"gopkg.in/Iwark/spreadsheet.v2"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var service *spreadsheet.Service

// получаем координаты из Яндекс справочника
// https://tech.yandex.ru/maps/geosearch/doc/concepts/request-docpage/
func GetCoordinatesFromYandex() {

	file, err := os.OpenFile("critic.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	spreadsheetID := config.SpreadsheetIDFsinPlaces
	sheet, err := service.FetchSpreadsheet(spreadsheetID)
	checkError(err)

	mainSheetFSIN, err := sheet.SheetByID(0)
	checkError(err)

	// 0 row это название полей, поэтому начинаем с 1 row
	// НЕ более 500 запросов в день к search-maps.yandex.ru
	for row := 1; row <= 882; row++ {
		var place model.Place

		place.Name = mainSheetFSIN.Rows[row][0].Value
		place.Type = mainSheetFSIN.Rows[row][1].Value
		place.Location = mainSheetFSIN.Rows[row][2].Value

		// https://tech.yandex.ru/maps/geosearch/doc/concepts/response_structure_business-docpage/
		url := "https://search-maps.yandex.ru/v1/"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
		}

		// Обязательными параметрами запроса являются: text, lang и apikey.
		q := req.URL.Query()
		q.Add("apikey", config.YandexAPIKey)
		q.Add("lang", "ru_RU")
		q.Add("text", place.Name+" "+place.Location)
		fmt.Println("делаем такой запрос:", place.Name+" "+place.Location)
		req.URL.RawQuery = q.Encode()

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		foundCount := gjson.Get(string(body), "properties.ResponseMetaData.SearchResponse.found").Uint()
		if foundCount > 1 {
			// записываем номер где возможны ошибки
			s := "предупреждение в записи номер: " + strconv.FormatUint(uint64(row), 10) + " \n"
			place.Warn = "требуется проверка!"
			_, err = file.WriteString(s)
			if err != nil {
				fmt.Println("возникла ошибка:", err)
			}
		} else if foundCount == 0 {
			place.Warn = "НЕ НАЙДЕНО!"
			s := "не найдено в записи номер: " + strconv.FormatUint(uint64(row), 10) + " \n"
			_, err = file.WriteString(s)
			if err != nil {
				fmt.Println("возникла ошибка:", err)
			}
			mainSheetFSIN.Update(row, 7, place.Warn)
			mainSheetFSIN.Update(row, 8, place.Warn)
			mainSheetFSIN.Update(row, 9, place.Warn)
			mainSheetFSIN.Update(row, 10, place.Warn)
			mainSheetFSIN.Update(row, 11, place.Warn)
			err = mainSheetFSIN.Synchronize()
			checkError(err)
			continue
		}

		fmt.Println("подходящих вариантов:", foundCount)

		// Контейнер результатов поиска. Обязательное поле.
		features := gjson.Get(string(body), "features").Array()

		// "coordinates":[
		// 132.337293, // [0] долгота
		// 43.987453 // [1] широта
		//]
		feature := features[0].Map()
		coordinates := feature["geometry"].Get("coordinates").Array()

		// долгота
		place.Position.Lng = coordinates[0].Float()
		// широта
		place.Position.Lat = coordinates[1].Float()

		companyMetaData := feature["properties"].Get("CompanyMetaData").Map()

		place.Address = companyMetaData["address"].String()
		place.Hours = companyMetaData["Hours"].Map()["text"].String()

		phones := companyMetaData["Phones"].Array()
		for ph := 0; ph < len(phones); ph++ {
			place.Phones = append(place.Phones, phones[ph].Map()["formatted"].String())
		}

		place.Website = companyMetaData["url"].String()

		mainSheetFSIN.Update(row, 7, strings.Join(place.Phones, ","))
		mainSheetFSIN.Update(row, 8, place.Hours)
		mainSheetFSIN.Update(row, 9, place.Website)
		mainSheetFSIN.Update(row, 10, place.Address)
		mainSheetFSIN.Update(row, 11, place.Warn)

		mainSheetFSIN.Update(row, 4, strconv.FormatFloat(place.Position.Lat, 'f', -1, 64))
		mainSheetFSIN.Update(row, 5, strconv.FormatFloat(place.Position.Lng, 'f', -1, 64))

		fmt.Printf("\n Place: %+v \n", place)
		err = mainSheetFSIN.Synchronize()
		checkError(err)
		time.Sleep(1 * time.Second)
	}
}

// убираем лишние скобки [] из Place.Notes
func UpdatePlaceNotes() {
	spreadsheetFsinPlaces := config.SpreadsheetIDFsinPlaces
	sheet, err := service.FetchSpreadsheet(spreadsheetFsinPlaces)
	checkError(err)
	fmt.Println("UpdatePlaceNotes...")
	sheetFSIN, err := sheet.SheetByID(0)
	checkError(err)
	for i := 123; i <= len(sheetFSIN.Rows)-1; i++ {
		Notes := sheetFSIN.Rows[i][3].Value
		Notes = regexp.MustCompile(`\\[.*?]`).ReplaceAllString(strings.Trim(Notes, "\n"), "")
		sheetFSIN.Update(i, 3, Notes)
		err = sheetFSIN.Synchronize()
		if err != nil {
			fmt.Println("ошибка на:", i)
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}
}

// получение координат по адресу by Google maps api
// метод не рекомендуется использовать из-за того что Goolge API часто не знает ФСИН учреждениe
// Deprecated:
func GetCoordinatesFromGoogle(address string) (float64, float64) {
	c, err := maps.NewClient(maps.WithAPIKey(config.GoogleMapsAPIKey))
	if err != nil {
		fmt.Printf("fatal error: %s", err)
	}
	geo := maps.GeocodingRequest{
		Address: address,
	}
	res, err := c.Geocode(context.TODO(), &geo)
	if err != nil {
		fmt.Println(err)
	}

	return res[0].Geometry.Location.Lat, res[0].Geometry.Location.Lng
}

// парсинг учреждений с wikipedia.org и добавление в Google Sheets
func WikiPlaces() {

	// Request the HTML page.
	res, err := http.Get("https://ru.wikipedia.org/wiki/%D0%A1%D0%BF%D0%B8%D1%81%D0%BE%D0%BA_%D0%BF%D0%B5%D0%BD%D0%B8%D1%82%D0%B5%D0%BD%D1%86%D0%B8%D0%B0%D1%80%D0%BD%D1%8B%D1%85_%D1%83%D1%87%D1%80%D0%B5%D0%B6%D0%B4%D0%B5%D0%BD%D0%B8%D0%B9_%D0%A0%D0%BE%D1%81%D1%81%D0%B8%D0%B8")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var places []model.Place
	var line = 1

	// Find the review items
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			var place model.Place

			rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
				switch indexth {
				case 0:
					place.Name = tablecell.Text()
				case 1:
					place.Type = tablecell.Text()
				case 2:
					place.Location = tablecell.Text()
				case 3:
					place.Notes = tablecell.Text()
				default:
					fmt.Println("default", tablecell.Text())
				}
			})
			if !strings.Contains(place.Name, "Название") {
				lat, lng := GetCoordinatesFromGoogle(place.Name)
				place.Position.Lng = lng
				place.Position.Lat = lat
				places = append(places, place)

				spreadsheetID := config.SpreadsheetIDFsinPlaces
				sheet, err := service.FetchSpreadsheet(spreadsheetID)
				checkError(err)

				mainSheetFSIN, err := sheet.SheetByID(0)
				checkError(err)

				mainSheetFSIN.Update(line, 0, place.Name)
				mainSheetFSIN.Update(line, 1, place.Type)
				mainSheetFSIN.Update(line, 2, place.Location)
				mainSheetFSIN.Update(line, 3, place.Notes)
				mainSheetFSIN.Update(line, 4, strconv.FormatFloat(place.Position.Lat, 'f', 6, 64))
				mainSheetFSIN.Update(line, 5, strconv.FormatFloat(place.Position.Lng, 'f', 6, 64))

				fmt.Println(place)
				time.Sleep(2 * time.Second)
				err = mainSheetFSIN.Synchronize()
				if err != nil {
					fmt.Println(err.Error())
				}
				line++
				fmt.Println("линия ", line)
			}
		})
	})
}

// подсчет кол-во нарушений для каждого ФСИН по нашим данным
func CountNumberOfViolations(c *gin.Context) {

	spreadsheetIDFsinPlaces := config.SpreadsheetIDForms
	sheet, err := service.FetchSpreadsheet(spreadsheetIDFsinPlaces)
	checkError(err)

	mainSheet, err := sheet.SheetByID(0)
	checkError(err)

	violations := make(map[string]uint64)
	for i := 1; i <= len(mainSheet.Rows)-1; i++ {
		val, ok := violations[mainSheet.Rows[i][32].Value+" "+mainSheet.Rows[i][33].Value]
		if ok {
			val++
			fmt.Println(mainSheet.Rows[i][32].Value + " " + mainSheet.Rows[i][33].Value)
			violations[mainSheet.Rows[i][32].Value+" "+mainSheet.Rows[i][33].Value] = val
		} else {
			violations[mainSheet.Rows[i][32].Value+" "+mainSheet.Rows[i][33].Value] = 1
		}
	}

	spreadsheetIDFsinPlaces = config.SpreadsheetIDFsinPlaces
	sheet, err = service.FetchSpreadsheet(spreadsheetIDFsinPlaces)
	checkError(err)

	mainSheet, err = sheet.SheetByID(0)
	checkError(err)

	for i := 1; i <= len(mainSheet.Rows)-1; i++ {
		val, ok := violations[mainSheet.Rows[i][4].Value+" "+mainSheet.Rows[i][5].Value]
		if ok {
			mainSheet.Update(i, 6, strconv.FormatUint(val, 10))
			fmt.Println(violations[mainSheet.Rows[i][4].Value+" "+mainSheet.Rows[i][5].Value])
			err := mainSheet.Synchronize()
			checkError(err)
			time.Sleep(1 * time.Second)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"violations": violations,
	})
}

// telephone and etc to table2 from Google API
// Deprecated:
func AddInfo() {
	c, err := maps.NewClient(maps.WithAPIKey(config.GoogleMapsAPIKey))
	if err != nil {
		fmt.Printf("fatal error: %s", err)
	}

	spreadsheetID := config.SpreadsheetIDFsinPlaces
	sheet, err := service.FetchSpreadsheet(spreadsheetID)
	checkError(err)

	mainSheet, err := sheet.SheetByID(0)
	checkError(err)

	for i := 1; i < len(mainSheet.Rows)-1; i++ {

		var r maps.FindPlaceFromTextRequest
		r.Input = mainSheet.Rows[i][0].Value
		r.InputType = "textquery"
		FindPlace, _ := c.FindPlaceFromText(context.TODO(), &r)

		if len(FindPlace.Candidates) == 0 {
			continue
		}

		var r2 maps.PlaceDetailsRequest
		r2.PlaceID = FindPlace.Candidates[0].PlaceID
		res, err := c.PlaceDetails(context.TODO(), &r2)
		checkError(err)
		fmt.Println(res.Website)

		mainSheet.Update(i, 7, res.FormattedPhoneNumber)
		mainSheet.Update(i, 8, strconv.FormatFloat(float64(res.Rating), 'f', 6, 64))
		mainSheet.Update(i, 9, res.Website)

		err = mainSheet.Synchronize()
		checkError(err)
		time.Sleep(1 * time.Second)
	}
}

// добавление координат в таблицу from Google API
// Google API неправильно находит координаты!
// Deprecated:
func AddCoordinatesToTable() {

	spreadsheetID := config.SpreadsheetIDForms
	sheet, err := service.FetchSpreadsheet(spreadsheetID)
	checkError(err)

	mainSheet, err := sheet.SheetByID(0)
	checkError(err)

	var fullName string

	for i := 1; i <= len(mainSheet.Rows)-1; i++ {
		region := mainSheet.Rows[i][2].Value
		FSINОrganization := mainSheet.Rows[i][3].Value
		if strings.ContainsRune(FSINОrganization, ',') {
			FSINОrganization = FSINОrganization[:strings.IndexByte(FSINОrganization, ',')]
		}
		fullName = region + " " + FSINОrganization

		fmt.Println("fullNAME:", fullName)

		lat, lng := GetCoordinatesFromGoogle(fullName)
		mainSheet.Update(i, 32, strconv.FormatFloat(lat, 'f', 6, 64))
		mainSheet.Update(i, 33, strconv.FormatFloat(lng, 'f', 6, 64))
		err := mainSheet.Synchronize()
		checkError(err)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

// места из изначальной таблицы
// Deprecated:
//func OldPlaces(c *gin.Context) {
//
//	places := make([]model.OldPlace, 0)
//
//	for i := 1; i <= len(mainSheet.Rows)-1; i++ {
//		var place model.OldPlace
//
//		place.Region = mainSheet.Rows[i][2].Value
//
//		place.FSINОrganization = mainSheet.Rows[i][3].Value
//
//		place.FullName = place.Region + " " + place.FSINОrganization
//
//		places = append(places, place)
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"places": places,
//	})
//}
