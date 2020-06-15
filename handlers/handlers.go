package handlers

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
	"strings"
	"time"
)

var Service *spreadsheet.Service

var mongoPlaces []model.Place
var mongoShortPlaces []model.ShortPlace
var mongoViolations []model.Violation
var mongoCoronaViolations []model.CoronaViolation

var violationsStats interface{}

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

var explanations = map[string]string{
	"total_count": "общее кол-во обращений",

	"physical_impact_from_employees": "С какими фактами применения физического воздействия со стороны сотрудников ФСИН Вам приходилось сталкиваться?",
	"physical_impact_from_prisoners": "С какими фактами применения физического воздействия со стороны заключенных Вам приходилось сталкиваться?",

	"psychological_impact_from_employees": "С какими фактами психологического воздействия со стороны сотрудников ФСИН Вам приходилось сталкиваться?",
	"psychological_impact_from_prisoners": "С какими фактами психологического воздействия со стороны заключенных Вам приходилось сталкиваться?",

	"can_prisoners_submit_complaints": "Есть ли у заключенных возможность направлять жалобы, ходатайства и заявления в надзирающие органы и правозащитные организации?",

	"communication_with_relatives": "Какие нарушения, связанные с иными формами общения с Родственниками, Вам известны?",
	"communication_with_lawyer":    "Какие нарушения, связанные с общением с адвокатом (иным лицом, имеющим право на оказание юридической помощи), Вам известны?",

	"visits_with_relatives": "Какие нарушения, связанные с предоставлением свиданий с Родственниками, Вам известны?",

	"corruption_from_employees": "Приходилось ли Вам сталкиваться с иными случаями коррупции сотрудников ФСИН?",
	"extortions_from_employees": "В каких случаях Вы сталкивались с фактами вымогательства со стороны сотрудников ФСИН?",
	"extortions_from_prisoners": "Приходилось ли Вам сталкиваться с фактами вымогательства со стороны заключенных?",

	"violations_of_medical_care": "Какие нарушения, связанные с оказанием медицинской помощи, Вы можете отметить?",
}

func Analytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"violations_stats":           violationsStats,
		"total_count_appeals":        db.CountAllViolations(),
		"total_count_appeals_corona": db.CountCoronaViolations(),
	})
}

func Explanations(c *gin.Context) {
	c.JSON(http.StatusOK, explanations)
}

// получение всех/одного ФСИН учреждений
func Places(c *gin.Context) {
	if c.Param("_id") == "" {
		c.JSON(http.StatusOK, gin.H{
			"places": mongoShortPlaces,
		})
	} else {
		for _, place := range mongoPlaces {
			if place.ID.Hex() == c.Param("_id") {
				if place.NumberOfViolations != 0 {
					for _, violation := range mongoViolations {
						for _, placeID := range violation.PlacesID {
							if placeID.Hex() == c.Param("_id") {
								place.Violations = append(place.Violations, violation)
							}
						}
					}
				}
				if place.Coronavirus {
					for _, corona := range mongoCoronaViolations {
						if corona.PlaceID == place.ID {
							place.CoronaViolations = append(place.CoronaViolations, corona)
						}
					}
				}
				c.JSON(http.StatusOK, gin.H{
					"place": place,
				})
				return
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "place not found",
		})
	}
}

// получение всех нарушений
func Violations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"violations": mongoViolations,
	})
}

func UpdateAllPlaces() {
	for {
		mongoShortPlaces = db.ShortPlaces()
		mongoPlaces = db.Places()
		mongoViolations = db.Violations()
		mongoCoronaViolations = db.CoronaViolations()
		violationsStats = db.CountViolations()

		fmt.Println("updated all, sleep for 5 minutes...")
		time.Sleep(5 * time.Minute)
	}
}

func CoronaPlaces(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"corona_violations": mongoCoronaViolations,
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
			if field.Type() != reflect.TypeOf("") || "add_files" == reflect.ValueOf(field).Type().Name() {
				continue
			}

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
