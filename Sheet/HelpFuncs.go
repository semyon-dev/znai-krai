// Здесь содержатся вспомогательные
// функции для анализа и сбора данных которые
// я написал для временного или единоразового использования

package Sheet

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/semyon-dev/RusSeated/Config"
	"github.com/semyon-dev/RusSeated/Model"
	"googlemaps.github.io/maps"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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

// подсчет кол-во нарушений для каждого ФСИН по нашим данным
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
		fmt.Println(res.Website)

		mainSheet.Update(i, 7, res.FormattedPhoneNumber)
		mainSheet.Update(i, 8, strconv.FormatFloat(float64(res.Rating), 'f', 6, 64))
		mainSheet.Update(i, 9, res.Website)

		err = mainSheet.Synchronize()
		checkError(err)
		time.Sleep(1 * time.Second)
	}
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

// места из изначальной таблицы Deprecated:
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
