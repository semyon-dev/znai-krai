package sheet

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/db"
	"github.com/semyon-dev/znai-krai/model"
	"go.mongodb.org/mongo-driver/bson"
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

// все места ФСИН учреждений
var places []model.Place
var placesBson []bson.M
var violations []bson.M

var placesCorona []model.PlaceCorona

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
	var form model.Form
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

// обновляем массив мест каждые 5 минут из Google Sheet всех учреждений
func UpdatePlaces() {
	spreadsheetFsinPlaces := config.SpreadsheetIDFsinPlaces
	sheet, err := Service.FetchSpreadsheet(spreadsheetFsinPlaces)
	checkError(err)
	fmt.Println("updating places...")
	sheetFSIN, err := sheet.SheetByID(0)
	checkError(err)
	places = nil
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

		for _, v := range placesCorona {
			if v.Position.Lat == place.Position.Lat && v.Position.Lng == place.Position.Lng {
				place.Coronavirus = true
			}
		}

		places = append(places, place)
	}
}

// обновляем массив мест вируса из Google Sheet всех учреждений
func UpdateCoronaPlaces() {
	spreadsheetFsinPlaces := config.SpreadsheetCoronavirus
	sheet, err := Service.FetchSpreadsheet(spreadsheetFsinPlaces)
	checkError(err)
	fmt.Println("updating corona places...")
	sheetFSIN, err := sheet.SheetByID(0)
	checkError(err)
	places = nil
	for i := 1; i <= len(sheetFSIN.Rows)-1; i++ {
		var place model.PlaceCorona

		place.Date = sheetFSIN.Rows[i][1].Value
		place.Info = sheetFSIN.Rows[i][4].Value
		place.CommentFSIN = sheetFSIN.Rows[i][5].Value

		place.Position.Lat, err = strconv.ParseFloat(sheetFSIN.Rows[i][10].Value, 64)
		if err != nil {
			place.Position.Lat = 0
		}
		place.Position.Lng, err = strconv.ParseFloat(sheetFSIN.Rows[i][11].Value, 64)
		if err != nil {
			place.Position.Lng = 0
		}
		placesCorona = append(placesCorona, place)
	}
}

func Analytics(c *gin.Context) {
	c.JSON(http.StatusOK, db.CountViolations())
}

// получение всех ФСИН учреждений
func Places(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"places": placesBson,
	})
}

func UpdateAllPlaces() {
	for {
		UpdateCoronaPlaces()
		UpdatePlaces()
		//db.UpdatePlaces(&places)
		placesBson = db.Places()
		violations = db.Violations()
		fmt.Println("updated all, sleep for 30 minutes...")
		time.Sleep(30 * time.Minute)
	}
}

// получение всех ФСИН учреждений
func CoronaPlaces(c *gin.Context) {
	if c.Query("lat") != "" && c.Query("lng") != "" {
		lat, err := strconv.ParseFloat(c.Query("lat"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			return
		}
		lng, err := strconv.ParseFloat(c.Query("lng"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			return
		}
		for _, v := range placesCorona {
			if v.Position.Lat == lat && v.Position.Lng == lng {
				c.JSON(http.StatusOK, gin.H{
					"places_corona": v,
				})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"places_corona": placesCorona,
	})
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
