// Здесь содержатся вспомогательные
// функции для анализа и сбора данных которые
// я написал для временного или единоразово использования
// В данным момент эти функции не используются и не поддерживаются, но могут понадобиться.
// Deprecated:
package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/db"
	"github.com/semyon-dev/znai-krai/model"
	"github.com/semyon-dev/znai-krai/sheet"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/mongo"
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

// deprecated:
var service *spreadsheet.Service

var sheetCoronaViolations []model.CoronaViolation

var sheetPlaces []model.Place

// обновляем нарушения коронавируса из Google Sheet всех учреждений
func UpdateCoronaPlaces() {
	spreadsheetFsinPlaces := config.SpreadsheetCoronavirus
	sheet, err := handlers.Service.FetchSpreadsheet(spreadsheetFsinPlaces)
	checkError(err)
	fmt.Println("updating corona sheetPlaces...")
	sheetCorona, err := sheet.SheetByID(0)
	checkError(err)
	sheetPlaces = nil
	for i := 1; i <= len(sheetCorona.Rows)-1; i++ {
		var coronaViolation model.CoronaViolation

		coronaViolation.Date = sheetCorona.Rows[i][1].Value
		coronaViolation.NameOfFSIN = sheetCorona.Rows[i][2].Value
		coronaViolation.Region = sheetCorona.Rows[i][3].Value
		coronaViolation.Info = sheetCorona.Rows[i][4].Value
		coronaViolation.CommentFSIN = sheetCorona.Rows[i][5].Value
		coronaViolation.Status = sheetCorona.Rows[i][12].Value

		coronaViolation.Position.Lat, err = strconv.ParseFloat(sheetCorona.Rows[i][10].Value, 64)
		if err != nil {
			coronaViolation.Position.Lat = 0
		}
		coronaViolation.Position.Lng, err = strconv.ParseFloat(sheetCorona.Rows[i][11].Value, 64)
		if err != nil {
			coronaViolation.Position.Lng = 0
		}
		sheetCoronaViolations = append(sheetCoronaViolations, coronaViolation)
	}
}

func UpdateCoronaPlacesToMongo() {
	spreadsheetCorona := config.SpreadsheetCoronavirus
	sheet, err := handlers.Service.FetchSpreadsheet(spreadsheetCorona)
	checkError(err)
	fmt.Println("updating corona sheetPlaces...")
	sheetCorona, err := sheet.SheetByID(0)
	checkError(err)
	sheetPlaces = nil
	for i := 1; i <= len(sheetCorona.Rows)-1; i++ {

		var coronaViolation model.CoronaViolation

		coronaViolation.Date = sheetCorona.Rows[i][1].Value
		coronaViolation.NameOfFSIN = sheetCorona.Rows[i][2].Value
		coronaViolation.Region = sheetCorona.Rows[i][3].Value
		coronaViolation.Info = sheetCorona.Rows[i][4].Value
		coronaViolation.CommentFSIN = sheetCorona.Rows[i][5].Value
		coronaViolation.Status = sheetCorona.Rows[i][12].Value

		coronaViolation.Position.Lat, err = strconv.ParseFloat(sheetCorona.Rows[i][10].Value, 64)
		if err != nil {
			coronaViolation.Position.Lat = 0
		}
		coronaViolation.Position.Lng, err = strconv.ParseFloat(sheetCorona.Rows[i][11].Value, 64)
		if err != nil {
			coronaViolation.Position.Lng = 0
		}
		for _, mongoPlace := range mongoPlaces {
			if coronaViolation.Position.Lat == mongoPlace.Position.Lat && coronaViolation.Position.Lng == mongoPlace.Position.Lng {
				coronaViolation.PlaceID = mongoPlace.ID
			}
		}
		db.AddCoronaViolation(coronaViolation)
	}
}

// обновляем массив мест из Google Sheet всех учреждений
func UpdateSheetPlaces() {
	spreadsheetFsinPlaces := config.SpreadsheetIDFsinPlaces
	sheet, err := handlers.Service.FetchSpreadsheet(spreadsheetFsinPlaces)
	checkError(err)
	fmt.Println("updating sheetPlaces...")
	sheetFSIN, err := sheet.SheetByID(0)
	checkError(err)
	sheetPlaces = nil
	for i := 1; i <= len(sheetFSIN.Rows)-1; i++ {
		var place model.Place

		place.Name = sheetFSIN.Rows[i][0].Value
		place.Type = sheetFSIN.Rows[i][1].Value
		place.Location = sheetFSIN.Rows[i][2].Value

		place.Notes = sheetFSIN.Rows[i][3].Value
		place.Notes = strings.Trim(place.Notes, "\n")

		place.Position.Lat, err = strconv.ParseFloat(sheetFSIN.Rows[i][4].Value, 64)
		if err != nil {
			place.Position.Lat = 0
		}
		place.Position.Lng, err = strconv.ParseFloat(sheetFSIN.Rows[i][5].Value, 64)
		if err != nil {
			place.Position.Lng = 0
		}

		place.NumberOfViolations, err = strconv.ParseUint(sheetFSIN.Rows[i][6].Value, 10, 64)
		if err != nil {
			place.NumberOfViolations = 0
		}

		place.Phones = strings.Split(sheetFSIN.Rows[i][7].Value, ",")
		place.Hours = sheetFSIN.Rows[i][8].Value
		place.Website = sheetFSIN.Rows[i][9].Value
		place.Address = sheetFSIN.Rows[i][10].Value
		place.Warn = sheetFSIN.Rows[i][11].Value

		for _, coronaViolation := range sheetCoronaViolations {
			if coronaViolation.Position.Lat == place.Position.Lat && coronaViolation.Position.Lng == place.Position.Lng {
				place.Coronavirus = true
			}
		}
		sheetPlaces = append(sheetPlaces, place)
	}
}

// обновляем нарушения в MongoDB из Google Sheet всех учреждений
func UpdateViolations() {
	spreadsheetFsinPlaces := config.SpreadsheetIDForms
	sheet, err := handlers.Service.FetchSpreadsheet(spreadsheetFsinPlaces)
	checkError(err)
	fmt.Println("updating violations...")
	sheetForms, err := sheet.SheetByID(0)
	checkError(err)
	for row := 1; row < 379; row++ {

		fmt.Println("--------", row, "--------")
		time.Sleep(200 * time.Millisecond)
		var form model.Violation
		form.Time = sheetForms.Rows[row][0].Value
		form.Status = sheetForms.Rows[row][1].Value
		form.Region = sheetForms.Rows[row][2].Value
		form.FSINOrganization = sheetForms.Rows[row][3].Value

		if strings.Contains(sheetForms.Rows[row][3].Value, ",") {
			fmt.Println("несколько МЛС")
			names := strings.Split(sheetForms.Rows[row][3].Value, ",")
			for _, name := range names {
				body := yandexRequest(form.Region + " " + name)
				foundCount := gjson.Get(string(body), "properties.ResponseMetaData.SearchResponse.found").Uint()
				switch foundCount {
				case 0:
					fmt.Println("не нашли")
					u := sheetForms.Rows[row][36].Value
					sheetForms.Update(row, 36, u+" "+"not found by yandex")
				case 1:
					fmt.Println("нашли однозначно")
					features := gjson.Get(string(body), "features").Array()
					feature := features[0].Map()
					coordinates := feature["geometry"].Get("coordinates").Array()

					var pos = model.Position{Lat: coordinates[1].Float(), Lng: coordinates[0].Float()}

					form.Positions = append(form.Positions, pos)

					place, err := db.FindPlace(pos)
					switch err {
					case mongo.ErrNoDocuments:
						fmt.Println("не найдено в MongoDB")
						u := sheetForms.Rows[row][36].Value
						sheetForms.Update(row, 36, u+","+"not found")
					case nil:
						form.PlacesID = append(form.PlacesID, place.ID)
					default:
						u := sheetForms.Rows[row][36].Value
						sheetForms.Update(row, 36, u+","+"mongo error")
					}
					if err != nil {
						err = sheetForms.Synchronize()
						checkErrorWithType("_sheetForms.Synchronize()", err)
					}
				default:
					features := gjson.Get(string(body), "features").Array()
					for i := 0; i < len(features); i++ {
						feature := features[i].Map()
						companyMetaData := feature["properties"].Get("CompanyMetaData").Map()
						fmt.Println("----------------------")
						fmt.Println("номер: ", i)
						fmt.Println("name: ", companyMetaData["name"].String())
						fmt.Println("address: ", companyMetaData["address"].String())
						fmt.Println("categories: ", companyMetaData["Categories"].Value())
						fmt.Println("description: ", companyMetaData["description"].String())
						fmt.Println("url: ", companyMetaData["url"].String())
					}
					fmt.Println("Какой добавить? (s - пропустить)")
					var choice string
					fmt.Scan(&choice)
					if choice != "s" {
						i, err := strconv.Atoi(choice)
						checkErrorWithType("strconv ", err)
						feature := features[i].Map()
						coordinates := feature["geometry"].Get("coordinates").Array()

						var pos = model.Position{Lat: coordinates[1].Float(), Lng: coordinates[0].Float()}

						form.Positions = append(form.Positions, pos)

						place, err := db.FindPlace(pos)
						switch err {
						case mongo.ErrNoDocuments:
							fmt.Println("не найдено в MongoDB")
							u := sheetForms.Rows[row][36].Value
							sheetForms.Update(row, 36, u+" "+"not found")
						case nil:
							form.PlacesID = append(form.PlacesID, place.ID)
						default:
							u := sheetForms.Rows[row][36].Value
							sheetForms.Update(row, 36, u+" "+"mongo error")
						}
						if err != nil {
							err = sheetForms.Synchronize()
							checkErrorWithType("_sheetForms.Synchronize()", err)
						}
					} else {
						continue
					}
				}
			}
			// position and place_id
			res, err := json.Marshal(form.Positions)
			fmt.Println(string(res))
			if err != nil {
				fmt.Println("JSON err ", err)
			}
			sheetForms.Update(row, 33, string(res))

			//res, err := primitive.ObjectIDFromHex(sheetForms.Rows[row][35].Value)
			//if err == nil || sheetForms.Rows[row][35].Value != "not found" {
			//	form.PlacesID = append(form.PlacesID, res)
			//}
			// position and place_id
			res, err = json.Marshal(form.PlacesID)
			if err != nil {
				fmt.Println("JSON err ", err)
			}
			sheetForms.Update(row, 35, string(res))

			err = sheetForms.Synchronize()
			checkErrorWithType("sheetForms.Synchronize()", err)

			form.Approved = true // у старых заявок стояло "Да"

			db.UpdateViolation(form)

			// если без запятой:
		} else {
			body := yandexRequest(form.Region + " " + form.FSINOrganization)
			foundCount := gjson.Get(string(body), "properties.ResponseMetaData.SearchResponse.found").Uint()
			if foundCount > 1 {
				features := gjson.Get(string(body), "features").Array()
				for i := 0; i < len(features); i++ {
					feature := features[i].Map()
					companyMetaData := feature["properties"].Get("CompanyMetaData").Map()
					fmt.Println("----------------------")
					fmt.Println("номер: ", i)
					fmt.Println("name: ", companyMetaData["name"].String())
					fmt.Println("address: ", companyMetaData["address"].String())
					fmt.Println("categories: ", companyMetaData["Categories"].Value())
					fmt.Println("description: ", companyMetaData["description"].String())
					fmt.Println("url: ", companyMetaData["url"].String())
				}
				fmt.Println("Какой добавить? (s - пропустить)")
				var choice string
				fmt.Scan(&choice)
				if choice != "s" {
					i, err := strconv.Atoi(choice)
					checkError(err)
					feature := features[i].Map()

					// "coordinates":[
					// 132.337293, // [0] долгота
					// 43.987453 // [1] широта
					//]
					//	feature := features[0].Map()
					coordinates := feature["geometry"].Get("coordinates").Array()

					var position model.Position
					position = model.Position{
						Lat: coordinates[1].Float(),
						Lng: coordinates[0].Float(),
					}
					// долгота
					form.Positions = append(form.Positions, position)

					sheetForms.Update(row, 32, strconv.FormatFloat(position.Lat, 'f', -1, 64))
					sheetForms.Update(row, 33, strconv.FormatFloat(position.Lng, 'f', -1, 64))
					//	sheetViolations.Update(row, 34, form.Warn)

					place, err := db.FindPlace(position)
					switch err {
					case mongo.ErrNoDocuments:
						fmt.Println("не найдено в MongoDB")
						sheetForms.Update(row, 35, "not found")
					case nil:
						sheetForms.Update(row, 35, place.ID.Hex())
					default:
						sheetForms.Update(row, 35, "mongo error")
					}
					err = sheetForms.Synchronize()
					checkError(err)
				} else {
					continue
				}
			} else if foundCount == 1 {
				fmt.Println("нашли однозначно")
				features := gjson.Get(string(body), "features").Array()
				feature := features[0].Map()
				coordinates := feature["geometry"].Get("coordinates").Array()

				var position model.Position
				position = model.Position{
					Lat: coordinates[1].Float(),
					Lng: coordinates[0].Float(),
				}
				// долгота
				form.Positions = append(form.Positions, position)

				sheetForms.Update(row, 32, strconv.FormatFloat(position.Lat, 'f', -1, 64))
				sheetForms.Update(row, 33, strconv.FormatFloat(position.Lng, 'f', -1, 64))

				place, err := db.FindPlace(position)
				switch err {
				case mongo.ErrNoDocuments:
					fmt.Println("не найдено в MongoDB")
					sheetForms.Update(row, 35, "not found")
				case nil:
					sheetForms.Update(row, 35, place.ID.Hex())
				default:
					sheetForms.Update(row, 35, "mongo error")
				}
				err = sheetForms.Synchronize()
				checkError(err)
			} else if foundCount == 0 {
				fmt.Println("не нашли...")
				sheetForms.Update(row, 35, "not found")
				sheetForms.Update(row, 34, "not found")
				continue
			}

			form.Approved = true // у старых заявок стояло "Да"

			db.UpdateViolation(form)
		}
	}
}

func yandexRequest(requestText string) (body []byte) {

	// https://tech.yandex.ru/maps/geosearch/doc/concepts/response_structure_business-docpage/
	myurl := "https://search-maps.yandex.ru/v1/"
	req, err := http.NewRequest("GET", myurl, nil)
	if err != nil {
		fmt.Println(err)
	}

	// Обязательными параметрами запроса являются: text, lang и apikey.
	q := req.URL.Query()
	q.Add("apikey", config.YandexAPIKey)
	q.Add("lang", "ru_RU")
	q.Add("text", requestText)
	fmt.Println(" ")
	fmt.Println("--------------------------------------")
	fmt.Println(" ")
	fmt.Println("делаем запрос:", requestText)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	return body
}

// РУЧНОЙ СПОСОБ - координаты для короны
func HandChooseCoordinatesFromYandexForCorona(row int) {

	spreadsheetID := config.SpreadsheetCoronavirus
	fetchSpreadsheet, err := handlers.Service.FetchSpreadsheet(spreadsheetID)
	checkError(err)

	mainSheetFSIN, err := fetchSpreadsheet.SheetByID(0)
	checkError(err)

	fmt.Println("номер (row): ", row)
	time.Sleep(222 * time.Millisecond)

	if mainSheetFSIN.Rows[row][12].Value != "требуется проверка!" {
		fmt.Println("!!! такие вот дела !!! :", mainSheetFSIN.Rows[row][12].Value)
	}

	var coronaViolation model.CoronaViolation

	coronaViolation.NameOfFSIN = mainSheetFSIN.Rows[row][2].Value
	coronaViolation.Region = mainSheetFSIN.Rows[row][3].Value

	body := yandexRequest(coronaViolation.Region + " " + coronaViolation.NameOfFSIN)
	foundCount := gjson.Get(string(body), "properties.ResponseMetaData.SearchResponse.found").Uint()
	if foundCount > 1 {
		fmt.Println("требуется проверка!")
		features := gjson.Get(string(body), "features").Array()

		for i := 0; i < len(features); i++ {
			feature := features[i].Map()
			companyMetaData := feature["properties"].Get("CompanyMetaData").Map()
			fmt.Println("----------------------")
			fmt.Println("номер: ", i)
			fmt.Println("name: ", companyMetaData["name"].String())
			fmt.Println("address: ", companyMetaData["address"].String())
			fmt.Println("categories: ", companyMetaData["Categories"].Value())
			fmt.Println("description: ", companyMetaData["description"].String())
			fmt.Println("url: ", companyMetaData["url"].String())
		}
		fmt.Println("Какой добавить? (s - пропустить)")
		var choice string
		_, err := fmt.Scan(&choice)
		checkError(err)
		if choice != "s" {
			i, err := strconv.Atoi(choice)
			checkError(err)
			feature := features[i].Map()

			coordinates := feature["geometry"].Get("coordinates").Array()

			// долгота
			coronaViolation.Position.Lng = coordinates[0].Float()
			// широта
			coronaViolation.Position.Lat = coordinates[1].Float()

			mainSheetFSIN.Update(row, 10, strconv.FormatFloat(coronaViolation.Position.Lat, 'f', -1, 64))
			mainSheetFSIN.Update(row, 11, strconv.FormatFloat(coronaViolation.Position.Lng, 'f', -1, 64))
			mainSheetFSIN.Update(row, 12, " ")

			err = mainSheetFSIN.Synchronize()
			checkError(err)
		} else {
			mainSheetFSIN.Update(row, 12, "ручная проверка!")
			err = mainSheetFSIN.Synchronize()
			checkError(err)
		}
	} else if foundCount == 1 {
		features := gjson.Get(string(body), "features").Array()
		feature := features[0].Map()

		coordinates := feature["geometry"].Get("coordinates").Array()

		// долгота
		coronaViolation.Position.Lng = coordinates[0].Float()
		// широта
		coronaViolation.Position.Lat = coordinates[1].Float()

		mainSheetFSIN.Update(row, 10, strconv.FormatFloat(coronaViolation.Position.Lat, 'f', -1, 64))
		mainSheetFSIN.Update(row, 11, strconv.FormatFloat(coronaViolation.Position.Lng, 'f', -1, 64))
		mainSheetFSIN.Update(row, 12, " ")

		err = mainSheetFSIN.Synchronize()
		checkError(err)

	} else if foundCount == 0 {
		fmt.Println("не нашли...")
		mainSheetFSIN.Update(row, 10, "ручной поиск!")
		mainSheetFSIN.Update(row, 11, "ручной поиск!")
		mainSheetFSIN.Update(row, 12, "ручной поиск!")
		err = mainSheetFSIN.Synchronize()
		checkError(err)
	}
}

// получаем координаты из Яндекс справочника для Коронавирусной таблицы
// https://tech.yandex.ru/maps/geosearch/doc/concepts/request-docpage/
func GetCoordinatesFromYandexForCoronavirus() {

	spreadsheetID := config.SpreadsheetCoronavirus
	mySheet, err := handlers.Service.FetchSpreadsheet(spreadsheetID)
	checkError(err)

	sheetCorona, err := mySheet.SheetByID(0)
	checkError(err)

	// нулевой row это название полей, поэтому начинаем с 1
	// НЕ более 500 запросов в день к search-maps.yandex.ru
	for row := 49; row <= 50; row++ {

		var coronaViolation model.CoronaViolation

		coronaViolation.NameOfFSIN = sheetCorona.Rows[row][2].Value
		coronaViolation.Region = sheetCorona.Rows[row][3].Value

		fmt.Println("row:", row)
		fmt.Println("CoronaViolation:", coronaViolation)

		body := yandexRequest(coronaViolation.Region + " " + coronaViolation.NameOfFSIN)
		foundCount := gjson.Get(string(body), "properties.ResponseMetaData.SearchResponse.found").Uint()
		fmt.Println("подходящих вариантов:", foundCount)

		if foundCount > 1 {
			HandChooseCoordinatesFromYandexForCorona(row)
		} else if foundCount == 0 {
			coronaViolation.Status = "НЕ НАЙДЕНО!"
			sheetCorona.Update(row, 10, coronaViolation.Status)
			sheetCorona.Update(row, 11, coronaViolation.Status)
			sheetCorona.Update(row, 12, coronaViolation.Status)
			err = sheetCorona.Synchronize()
			checkError(err)
			continue
		}

		// Контейнер результатов поиска. Обязательное поле.
		features := gjson.Get(string(body), "features").Array()

		feature := features[0].Map()
		coordinates := feature["geometry"].Get("coordinates").Array()

		// долгота
		coronaViolation.Position.Lng = coordinates[0].Float()
		// широта
		coronaViolation.Position.Lat = coordinates[1].Float()

		sheetCorona.Update(row, 10, strconv.FormatFloat(coronaViolation.Position.Lat, 'f', -1, 64))
		sheetCorona.Update(row, 11, strconv.FormatFloat(coronaViolation.Position.Lng, 'f', -1, 64))

		fmt.Printf("\n Place: %+v \n", coronaViolation)
		err = sheetCorona.Synchronize()
		checkError(err)
		time.Sleep(11 * time.Millisecond)
	}
}

// РУЧНОЙ СПОСОБ - выбираем координаты из Яндекс справочника
// https://tech.yandex.ru/maps/geosearch/doc/concepts/request-docpage/
func ChooseCoordinatesFromYandex() {

	spreadsheetID := config.SpreadsheetIDFsinPlaces
	fetchSpreadsheet, err := handlers.Service.FetchSpreadsheet(spreadsheetID)
	checkError(err)

	mainSheetFSIN, err := fetchSpreadsheet.SheetByID(0)
	checkError(err)

	// 0 row это название полей, поэтому начинаем с 1 row
	// НЕ более 500 запросов в день к search-maps.yandex.ru
	for row := 1; row <= 882; row++ {

		fmt.Println("номер (row): ", row)
		time.Sleep(555 * time.Millisecond)

		if mainSheetFSIN.Rows[row][11].Value == "" {
			continue
		}
		if mainSheetFSIN.Rows[row][11].Value != "требуется проверка!" {
			fmt.Println("!!! такие вот дела !!! :", mainSheetFSIN.Rows[row][11].Value)
		}

		var place model.Place

		place.Name = mainSheetFSIN.Rows[row][0].Value
		place.Type = mainSheetFSIN.Rows[row][1].Value
		place.Location = mainSheetFSIN.Rows[row][2].Value

		// https://tech.yandex.ru/maps/geosearch/doc/concepts/response_structure_business-docpage/
		myurl := "https://search-maps.yandex.ru/v1/"
		req, err := http.NewRequest("GET", myurl, nil)
		if err != nil {
			fmt.Println(err)
		}

		// Обязательными параметрами запроса являются: text, lang и apikey.
		q := req.URL.Query()
		q.Add("apikey", config.YandexAPIKey)
		q.Add("lang", "ru_RU")
		q.Add("text", place.Name+" "+place.Location)
		fmt.Println(" ")
		fmt.Println("--------------------------------------")
		fmt.Println(" ")
		fmt.Println("делаем запрос:", place.Name+" "+place.Location)
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
		//	fmt.Println("подходящих вариантов:", foundCount)

		if foundCount > 1 {
			//	fmt.Println("требуется проверка!")
			// Контейнер результатов поиска. Обязательное поле.
			features := gjson.Get(string(body), "features").Array()

			for i := 0; i < len(features); i++ {
				feature := features[i].Map()
				companyMetaData := feature["properties"].Get("CompanyMetaData").Map()
				fmt.Println("----------------------")
				fmt.Println("номер: ", i)
				fmt.Println("name: ", companyMetaData["name"].String())
				fmt.Println("address: ", companyMetaData["address"].String())
				fmt.Println("categories: ", companyMetaData["Categories"].Value())
				fmt.Println("description: ", companyMetaData["description"].String())
				fmt.Println("url: ", companyMetaData["url"].String())
			}
			fmt.Println("Какой добавить? (s - пропустить)")
			var choice string
			fmt.Scan(&choice)
			if choice != "s" {
				i, err := strconv.Atoi(choice)
				checkError(err)
				feature := features[i].Map()
				companyMetaData := feature["properties"].Get("CompanyMetaData").Map()

				place.Address = companyMetaData["address"].String()
				place.Hours = companyMetaData["Hours"].Map()["text"].String()

				phones := companyMetaData["Phones"].Array()
				for ph := 0; ph < len(phones); ph++ {
					place.Phones = append(place.Phones, phones[ph].Map()["formatted"].String())
				}

				place.Website = companyMetaData["url"].String()

				// "coordinates":[
				// 132.337293, // [0] долгота
				// 43.987453 // [1] широта
				//]
				//	feature := features[0].Map()
				coordinates := feature["geometry"].Get("coordinates").Array()

				// долгота
				place.Position.Lng = coordinates[0].Float()
				// широта
				place.Position.Lat = coordinates[1].Float()

				mainSheetFSIN.Update(row, 7, strings.Join(place.Phones, ","))
				mainSheetFSIN.Update(row, 8, place.Hours)
				mainSheetFSIN.Update(row, 9, place.Website)
				mainSheetFSIN.Update(row, 10, place.Address)
				mainSheetFSIN.Update(row, 11, place.Warn)

				mainSheetFSIN.Update(row, 4, strconv.FormatFloat(place.Position.Lat, 'f', -1, 64))
				mainSheetFSIN.Update(row, 5, strconv.FormatFloat(place.Position.Lng, 'f', -1, 64))

				err = mainSheetFSIN.Synchronize()
				checkError(err)
			} else {
				continue
			}
		} else if foundCount == 0 {
			fmt.Println("не нашли...")
			continue
		} else {
			continue
		}
	}
}

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

// РУЧНОЙ СПОСОБ - выбираем координаты из Яндекс справочника для таблицы нарушений
// https://tech.yandex.ru/maps/geosearch/doc/concepts/request-docpage/
// Deprecated:
func ChooseCoordinatesFromYandexForViolations() {

	spreadsheetID := config.SpreadsheetIDForms
	fetchSpreadsheet, err := handlers.Service.FetchSpreadsheet(spreadsheetID)
	checkError(err)

	sheetViolations, err := fetchSpreadsheet.SheetByID(0)
	checkError(err)

	// 0 row это название полей, поэтому начинаем с 1 row
	// НЕ более 500 запросов в день к search-maps.yandex.ru
	// TODO: от 115 несколько учреждений раз раз сразу
	for row := 115; row < 116; row++ {

		fmt.Println("номер (row): ", row)
		time.Sleep(444 * time.Millisecond)

		var form model.Violation

		form.Region = sheetViolations.Rows[row][2].Value
		form.FSINOrganization = sheetViolations.Rows[row][3].Value

		body := yandexRequest(form.Region + " " + form.FSINOrganization)

		foundCount := gjson.Get(string(body), "properties.ResponseMetaData.SearchResponse.found").Uint()
		//	fmt.Println("подходящих вариантов:", foundCount)

		if foundCount > 1 {
			//	fmt.Println("требуется проверка!")
			// Контейнер результатов поиска. Обязательное поле.
			features := gjson.Get(string(body), "features").Array()

			for i := 0; i < len(features); i++ {
				feature := features[i].Map()
				companyMetaData := feature["properties"].Get("CompanyMetaData").Map()
				fmt.Println("----------------------")
				fmt.Println("номер: ", i)
				fmt.Println("name: ", companyMetaData["name"].String())
				fmt.Println("address: ", companyMetaData["address"].String())
				fmt.Println("categories: ", companyMetaData["Categories"].Value())
				fmt.Println("description: ", companyMetaData["description"].String())
				fmt.Println("url: ", companyMetaData["url"].String())
			}
			fmt.Println("Какой добавить? (s - пропустить)")
			var choice string
			fmt.Scan(&choice)
			if choice != "s" {
				i, err := strconv.Atoi(choice)
				checkError(err)
				feature := features[i].Map()

				// "coordinates":[
				// 132.337293, // [0] долгота
				// 43.987453 // [1] широта
				//]
				//	feature := features[0].Map()
				coordinates := feature["geometry"].Get("coordinates").Array()

				var position model.Position
				position = model.Position{
					Lat: coordinates[1].Float(),
					Lng: coordinates[0].Float(),
				}
				// долгота
				form.Positions = append(form.Positions, position)

				sheetViolations.Update(row, 32, strconv.FormatFloat(form.Positions[0].Lat, 'f', -1, 64))
				sheetViolations.Update(row, 33, strconv.FormatFloat(form.Positions[0].Lng, 'f', -1, 64))
				//	sheetViolations.Update(row, 34, form.Warn)

				place, err := db.FindPlace(position)
				switch err {
				case mongo.ErrNoDocuments:
					fmt.Println("не найдено в MongoDB")
					sheetViolations.Update(row, 35, "not found")
				case nil:
					sheetViolations.Update(row, 35, place.ID.Hex())
				default:
					sheetViolations.Update(row, 35, "mongo error")
				}
				err = sheetViolations.Synchronize()
				checkError(err)
			} else {
				continue
			}
		} else if foundCount == 1 {
			fmt.Println("нашли однозначно")
			features := gjson.Get(string(body), "features").Array()
			feature := features[0].Map()
			coordinates := feature["geometry"].Get("coordinates").Array()

			var position model.Position
			position = model.Position{
				Lat: coordinates[1].Float(),
				Lng: coordinates[0].Float(),
			}
			// долгота
			form.Positions = append(form.Positions, position)

			sheetViolations.Update(row, 32, strconv.FormatFloat(form.Positions[0].Lat, 'f', -1, 64))
			sheetViolations.Update(row, 33, strconv.FormatFloat(form.Positions[0].Lng, 'f', -1, 64))

			place, err := db.FindPlace(position)
			switch err {
			case mongo.ErrNoDocuments:
				fmt.Println("не найдено в MongoDB")
				sheetViolations.Update(row, 35, "not found")
			case nil:
				sheetViolations.Update(row, 35, place.ID.Hex())
			default:
				sheetViolations.Update(row, 35, "mongo error")
			}
			err = sheetViolations.Synchronize()
			checkError(err)
		} else if foundCount == 0 {
			fmt.Println("не нашли...")
			sheetViolations.Update(row, 35, "not found")
			sheetViolations.Update(row, 34, "not found")
			continue
		}
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

func checkErrorWithType(s string, err error) {
	if err != nil {
		fmt.Println("Error: ", s, err.Error())
	}
}
