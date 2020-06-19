package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/db"
	"github.com/semyon-dev/znai-krai/log"
	"github.com/semyon-dev/znai-krai/model"
	"github.com/semyon-dev/znai-krai/sheet"
	"googlemaps.github.io/maps"
	"net/http"
	"time"
)

var mongoPlaces []model.Place
var mongoShortPlaces []model.ShortPlace
var mongoViolations []model.Violation
var mongoCoronaViolations []model.CoronaViolation

var violationsStats interface{}

func Analytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"violations_stats":           violationsStats,
		"total_count_appeals":        db.CountAllViolations(),
		"total_count_appeals_corona": db.CountCoronaViolations(),
	})
}

func Explanations(c *gin.Context) {
	c.JSON(http.StatusOK, model.Explanations)
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
func Reviews(c *gin.Context) {

	cMaps, err := maps.NewClient(maps.WithAPIKey(config.GoogleMapsAPIKey))
	if err != nil {
		log.HandleErr(err)
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
func NewForm(c *gin.Context) {
	var form model.Violation
	var message string
	var status int
	err := c.ShouldBind(&form)
	if err != nil {
		status = http.StatusBadRequest
		message = "error: " + err.Error()
		checkError(err)
	} else {
		form.Source = "сайт"
		err = sheet.AddViolation(form)
		if err == nil {
			status = http.StatusOK
			message = "ok"
		} else {
			status = http.StatusInternalServerError
			message = "error " + err.Error()
		}
	}
	c.JSON(status, gin.H{
		"message": message,
	})
}

// новая форма нарушения по коронавируса
func NewFormCorona(c *gin.Context) {
	var form model.CoronaViolation
	var message string
	var status int
	err := c.ShouldBindJSON(&form)
	if err != nil {
		status = http.StatusBadRequest
		message = "error: " + err.Error()
		log.HandleErr(err)
	} else {
		form.Source = "сайт"
		err = sheet.AddCoronaViolation(form)
		if err == nil {
			status = http.StatusOK
			message = "ok"
		} else {
			status = http.StatusInternalServerError
			message = "InternalServerError"
		}
	}
	c.JSON(status, gin.H{
		"message": message,
	})
}

// новая форма нарушения по коронавируса
func NewMailing(c *gin.Context) {
	var form model.Mailing
	var message string
	var status int
	err := c.ShouldBindJSON(&form)
	if err != nil {
		status = http.StatusBadRequest
		message = "error: " + err.Error()
		log.HandleErr(err)
	} else {
		err = sheet.AddMailing(form)
		if err == nil {
			status = http.StatusOK
			message = "ok"
		} else {
			status = http.StatusInternalServerError
			message = "InternalServerError"
		}
	}
	c.JSON(status, gin.H{
		"message": message,
	})
}

func checkError(err error) {
	if err != nil {
		log.HandleErr(err)
	}
}
