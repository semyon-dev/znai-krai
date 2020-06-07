package sheet

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/db"
	"github.com/semyon-dev/znai-krai/model"
	"golang.org/x/oauth2/google"
	"googlemaps.github.io/maps"
	"gopkg.in/Iwark/spreadsheet.v2"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var Service *spreadsheet.Service

var sheetPlaces []model.Place
var sheetCoronaViolations []model.CoronaViolation

var mongoPlaces []model.Place
var mongoViolations []model.Violation
var mongoCoronaViolations []model.CoronaViolation

// Connect to Google Sheets
func Connect() {
	var data []byte
	var err error
	// если нет переменных окружения значит читаем файл credentials.json
	if len(config.Credentials) == 0 {
		data, err = ioutil.ReadFile("credentials.json")
		fmt.Println("read credentials from file")
		if err != nil {
			fmt.Println("критическая ошибка: не удалось импортировать переменные:", err)
		}
	} else {
		var f model.CredentialsFile
		fmt.Println("read credentials from env var")
		f.Type = "service_account"
		f.ProjectID = "zekovnet"
		f.PrivateKeyID = os.Getenv("private_key_id")
		f.PrivateKey = strings.ReplaceAll(os.Getenv("private_key"), "\\n", "\n")
		f.ClientEmail = os.Getenv("client_email")
		f.ClientID = os.Getenv("client_id")
		f.TokenURL = os.Getenv("token_uri")
		data, err = json.Marshal(f)
		checkError(err)
	}

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)

	client := conf.Client(context.TODO())
	Service = spreadsheet.NewServiceWithClient(client)

	formSpreadsheet, err := Service.FetchSpreadsheet(config.SpreadsheetIDForms)
	checkError(err)
	spreadsheetCoronavirus, err := Service.FetchSpreadsheet(config.SpreadsheetCoronavirus)
	checkError(err)
	fsinPlacesSpreadsheet, err := Service.FetchSpreadsheet(config.SpreadsheetIDFsinPlaces)
	checkError(err)
	fmt.Println("таблица нарушений (форм):", formSpreadsheet.Properties.Title)
	fmt.Println("таблица ФСИН учреждений:", fsinPlacesSpreadsheet.Properties.Title)
	fmt.Println("таблица с информацией по Коронавирусу:", spreadsheetCoronavirus.Properties.Title)
}

// обновляем массив мест из Google Sheet всех учреждений
func UpdateSheetPlaces() {
	spreadsheetFsinPlaces := config.SpreadsheetIDFsinPlaces
	sheet, err := Service.FetchSpreadsheet(spreadsheetFsinPlaces)
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

// обновляем нарушения коронавируса из Google Sheet всех учреждений
func UpdateCoronaPlaces() {
	spreadsheetFsinPlaces := config.SpreadsheetCoronavirus
	sheet, err := Service.FetchSpreadsheet(spreadsheetFsinPlaces)
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
	sheet, err := Service.FetchSpreadsheet(spreadsheetCorona)
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

func Analytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"violations_stats": db.CountViolations(),
		"total_count":      db.CountAllViolations(),
	})
}

// получение всех/одного ФСИН учреждений
func Places(c *gin.Context) {
	if c.Param("_id") == "" {
		c.JSON(http.StatusOK, gin.H{
			"places": mongoPlaces,
		})
	} else {
		for _, v := range mongoPlaces {
			if v.ID.Hex() == c.Param("_id") {
				for _, violation := range mongoViolations {
					for _, placeID := range violation.PlacesID {
						if placeID.Hex() == c.Param("_id") {
							v.Violations = append(v.Violations, violation)
						}
					}
				}
				c.JSON(http.StatusOK, gin.H{
					"place": v,
				})
				return
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "place not found",
			"place":   "",
		})
	}
}

// получение всех нарушений
func Violations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"mongoViolations": mongoViolations,
	})
}

func UpdateAllPlaces() {
	for {
		// UpdateCoronaPlaces()
		// UpdateSheetPlaces()
		//db.UpdateSheetPlaces(&sheetPlaces)

		mongoPlaces = db.Places()
		mongoViolations = db.Violations()
		mongoCoronaViolations = db.CoronaViolations()

		//for _, mongoPlace := range mongoPlaces {
		//	for _, violation := range mongoViolations {
		//		for _, placeID := range violation.PlacesID {
		//			if placeID.Hex() == mongoPlace.ID.Hex() {
		//				mongoPlace.NumberOfViolations++
		//			}
		//		}
		//	}
		//	db.UpdatePlace(mongoPlace)
		//}

		fmt.Println("updated all, sleep for 10 minutes...")
		time.Sleep(10 * time.Minute)
	}
}

// получение всех ФСИН учреждений
func CoronaPlaces(c *gin.Context) {
	if c.Query("lat") != "" && c.Query("lng") != "" {
		lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
		lng, err2 := strconv.ParseFloat(c.Query("lng"), 64)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			return
		}
		for _, v := range sheetCoronaViolations {
			if v.Position.Lat == lat && v.Position.Lng == lng {
				c.JSON(http.StatusOK, gin.H{
					"places_corona": v,
				})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"places_corona": mongoCoronaViolations,
	})
}

// получение отзывов с Google Maps
// TODO: refactor & testing
func Reviews(c *gin.Context) {

	cMaps, err := maps.NewClient(maps.WithAPIKey(config.GoogleMapsAPIKey))
	if err != nil {
		fmt.Printf("fatal error: %s", err)
	}

	var r maps.FindPlaceFromTextRequest
	r.Input = c.Param("name")
	r.InputType = "textquery"
	FindPlace, _ := cMaps.FindPlaceFromText(context.TODO(), &r)

	if len(FindPlace.Candidates) == 0 || len(FindPlace.Candidates) > 1 {
		c.JSON(http.StatusOK, gin.H{
			"reviews": nil,
		})
	} else {
		var r2 maps.PlaceDetailsRequest
		r2.PlaceID = FindPlace.Candidates[0].PlaceID
		res, err := cMaps.PlaceDetails(context.TODO(), &r2)
		checkError(err)
		c.JSON(http.StatusOK, gin.H{
			"reviews": res.Reviews,
		})
	}
}

// новая форма нарушения
// TODO: refactor & testing
func NewForm(c *gin.Context) {
	var form model.Violation
	var message string
	var status int
	err := c.ShouldBind(&form)

	if err != nil {
		fmt.Println(err.Error())
		status = 400
		message = "error: " + err.Error()
	} else {
		form.Source = "сайт"

		// We need a pointer so that we can set the value via reflection
		msValuePtr := reflect.ValueOf(&form)
		msValue := msValuePtr.Elem()

		// нужно для синхронизации
		formSpreadsheet, err := Service.FetchSpreadsheet(config.SpreadsheetIDForms)
		checkError(err)

		formsSheet, err := formSpreadsheet.SheetByID(0)
		checkError(err)

		row := len(formsSheet.Rows)
		column := 0
		currentTime := time.Now()

		formsSheet.Update(row, column, currentTime.Format("2006.01.02 15:04:05"))
		column++

		for ; column < msValue.NumField(); column++ {
			field := msValue.Field(column)

			// Ignore fields that don't have the same type as a string
			//if field.Type() != reflect.TypeOf("") {
			//	continue
			//}

			str := field.Interface().(string)
			str = strings.TrimSpace(str)
			field.SetString(str)

			// добавляем в таблицу
			formsSheet.Update(row, column, field.String())
		}

		err = formsSheet.Synchronize()
		if err == nil {
			status = 200
			message = "ok"
		} else {
			status = 400
			fmt.Println(err.Error())
			message = "error " + err.Error()
		}
	}
	c.JSON(status, gin.H{
		"message": message,
	})
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error ", err.Error())
	}
}

func checkErrorWithType(_type string, err error) {
	if err != nil {
		fmt.Println(_type, err.Error())
	}
}
