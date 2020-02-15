package Sheet

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/semyon-dev/RusSeated/Config"
	"github.com/semyon-dev/RusSeated/Model"
	"golang.org/x/oauth2/google"
	"googlemaps.github.io/maps"
	"gopkg.in/Iwark/spreadsheet.v2"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var mainSheet *spreadsheet.Sheet
var sheet spreadsheet.Spreadsheet
var service *spreadsheet.Service
var spreadsheetID string

func Connect() {
	data, err := ioutil.ReadFile("credentials.json")
	checkError(err)

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)

	client := conf.Client(context.TODO())
	service = spreadsheet.NewServiceWithClient(client)

	spreadsheetID = Config.SpreadsheetID
	sheet, err = service.FetchSpreadsheet(spreadsheetID)

	mainSheet, err = sheet.SheetByID(0)
	checkError(err)

	fmt.Println("sheet id or error:", sheet.ID)
}

// telephone and etc to table2
func AddInfo() {
	c, err := maps.NewClient(maps.WithAPIKey(Config.GoogleMapsAPIKey))
	if err != nil {
		fmt.Printf("fatal error: %s", err)
	}

	spreadsheetID = Config.SpreadsheetID_FSINPlaces
	sheet, err = service.FetchSpreadsheet(spreadsheetID)

	mainSheet, err = sheet.SheetByID(0)

	for i := 1; i < len(mainSheet.Rows)-1; i++ {

		fmt.Println(i)
		var r maps.FindPlaceFromTextRequest
		//r.Input = mainSheet.Rows[i][4].Value + " " + mainSheet.Rows[i][5].Value
		r.Input = mainSheet.Rows[i][0].Value
		fmt.Println(r.Input)
		r.InputType = "textquery"
		FindPlace, _ := c.FindPlaceFromText(context.TODO(), &r)
		//fmt.Println(FindPlace.Candidates[0].PlaceID)

		if len(FindPlace.Candidates) == 0{
			continue
		}

		var r2 maps.PlaceDetailsRequest
		r2.PlaceID = FindPlace.Candidates[0].PlaceID
		res, err := c.PlaceDetails(context.TODO(), &r2)
		fmt.Println(res.Website)

		mainSheet.Update(i, 7, res.FormattedPhoneNumber)
		mainSheet.Update(i, 8, strconv.FormatFloat(float64(res.Rating), 'f', 6, 64))
		mainSheet.Update(i, 9, res.Website)

		err = mainSheet.Synchronize()
		checkError(err)
		time.Sleep(1 * time.Second)
	}
}

func CountNumberOfViolations(c *gin.Context) {

	spreadsheetID = Config.SpreadsheetID
	sheet, err := service.FetchSpreadsheet(spreadsheetID)
	checkError(err)

	mainSheet, err = sheet.SheetByID(0)

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

	spreadsheetID = Config.SpreadsheetID_FSINPlaces
	sheet, err = service.FetchSpreadsheet(spreadsheetID)
	checkError(err)

	mainSheet, err = sheet.SheetByID(0)

	fmt.Println(violations["55.052166 83.004823"])

	for i := 1; i <= len(mainSheet.Rows)-1; i++ {
		//fmt.Println(i)
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

// добавление координат в таблицу
func AddCoordinatesToTable() {

	spreadsheetID = Config.SpreadsheetID
	sheet, err := service.FetchSpreadsheet(spreadsheetID)

	mainSheet, err := sheet.SheetByID(0)
	checkError(err)

	var fullName string

	for i := 1; i <= len(mainSheet.Rows)-1; i++ {
		fullName = ""
		region := mainSheet.Rows[i][2].Value
		FSINОrganization := mainSheet.Rows[i][3].Value
		if strings.ContainsRune(FSINОrganization, ',') {
			FSINОrganization = FSINОrganization[:strings.IndexByte(FSINОrganization, ',')]
		}
		fullName = region + " " + FSINОrganization

		fmt.Println("fullNAME:", fullName)

		lat, lng := GetCoordinates(fullName)
		mainSheet.Update(i, 32, strconv.FormatFloat(lat, 'f', 6, 64))
		mainSheet.Update(i, 33, strconv.FormatFloat(lng, 'f', 6, 64))
		err := mainSheet.Synchronize()
		checkError(err)
	}
}

// получение координат по адресу
func GetCoordinates(address string) (float64, float64) {
	c, err := maps.NewClient(maps.WithAPIKey(Config.GoogleMapsAPIKey))
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

// новая форма нарушения
func NewForm(c *gin.Context) {
	var form Model.Form
	var message string
	err := c.ShouldBindJSON(&form)
	if err != nil {
		fmt.Println(err.Error())
		message = "error: " + err.Error()
	} else {
		// We need a pointer so that we can set the value via reflection
		msValuePtr := reflect.ValueOf(&form)
		msValue := msValuePtr.Elem()

		// нужно для синхронизации
		sheet, err = service.FetchSpreadsheet(spreadsheetID)
		checkError(err)

		mainSheet, err = sheet.SheetByID(0)
		checkError(err)

		row := len(mainSheet.Rows)
		column := 0
		mainSheet.Update(row, column, strconv.FormatInt(int64(time.Now().Year()), 10))
		column++

		for ; column < msValue.NumField(); column++ {
			field := msValue.Field(column)

			// Ignore fields that don't have the same type as a string
			if field.Type() != reflect.TypeOf("") {
				continue
			}

			str := field.Interface().(string)
			str = strings.TrimSpace(str)
			field.SetString(str)

			// добавляем в таблицу
			mainSheet.Update(row, column, field.String())
		}

		err = mainSheet.Synchronize()
		if err == nil {
			message = "ok"
		} else {
			fmt.Println(err.Error())
			message = "error " + err.Error()
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

// получение всех ФСИН учреждений
func Places(c *gin.Context) {

	spreadsheetID = Config.SpreadsheetID_FSINPlaces
	sheet, err := service.FetchSpreadsheet(spreadsheetID)

	mainSheetFSIN, err := sheet.SheetByID(0)
	checkError(err)

	places := make([]Model.Place, 0)

	for i := 1; i <= 881; i++ {
		var place Model.Place

		place.Name = mainSheetFSIN.Rows[i][0].Value
		place.Type = mainSheetFSIN.Rows[i][1].Value
		place.Location = mainSheetFSIN.Rows[i][2].Value

		place.Notes = mainSheetFSIN.Rows[i][3].Value
		place.Notes = strings.Trim(place.Notes, "\n")

		place.Position.Lat, err = strconv.ParseFloat(mainSheetFSIN.Rows[i][4].Value, 64)
		place.Position.Lng, err = strconv.ParseFloat(mainSheetFSIN.Rows[i][5].Value, 64)

		place.NumberOfViolations, err = strconv.ParseUint(mainSheetFSIN.Rows[i][6].Value, 10, 64)

		place.PhoneNumber = mainSheetFSIN.Rows[i][7].Value
		place.GoogleMapsRating, err = strconv.ParseFloat(mainSheetFSIN.Rows[i][8].Value, 64)
		place.Website = mainSheetFSIN.Rows[i][9].Value

		places = append(places, place)
	}
	c.JSON(http.StatusOK, gin.H{
		"places": places,
	})
}

// места из изначальной таблицы
func OldPlaces(c *gin.Context) {

	places := make([]Model.OldPlace, 0)

	for i := 1; i <= len(mainSheet.Rows)-1; i++ {
		var place Model.OldPlace

		place.Region = mainSheet.Rows[i][2].Value

		place.FSINОrganization = mainSheet.Rows[i][3].Value

		place.FullName = place.Region + " " + place.FSINОrganization

		places = append(places, place)
	}
	c.JSON(http.StatusOK, gin.H{
		"places": places,
	})
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

// парсинг учреждений с википедии и добавление в Google Sheets
func WikiPlaces() {

	// Request the HTML page.
	res, err := http.Get("https://ru.wikipedia.org/wiki/%D0%A1%D0%BF%D0%B8%D1%81%D0%BE%D0%BA_%D0%BF%D0%B5%D0%BD%D0%B8%D1%82%D0%B5%D0%BD%D1%86%D0%B8%D0%B0%D1%80%D0%BD%D1%8B%D1%85_%D1%83%D1%87%D1%80%D0%B5%D0%B6%D0%B4%D0%B5%D0%BD%D0%B8%D0%B9_%D0%A0%D0%BE%D1%81%D1%81%D0%B8%D0%B8")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var places []Model.Place
	var line = 1

	// Find the review items
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			var place Model.Place

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
				lat, lng := GetCoordinates(place.Name)
				place.Position.Lng = lng
				place.Position.Lat = lat
				places = append(places, place)

				spreadsheetID = Config.SpreadsheetID_FSINPlaces
				sheet, err = service.FetchSpreadsheet(spreadsheetID)

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
