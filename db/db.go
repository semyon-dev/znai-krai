package db

import (
	"context"
	"fmt"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"strings"
	"time"
)

var db *mongo.Database

func Connect() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://"+config.MongoDBLogin+":"+config.MongoDBPass+"@main-h6nko.mongodb.net/test?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}

	// Create connect
	err = client.Connect(context.Background())
	if err != nil {
		fmt.Println("client MongoDB:", err)
	} else {
		fmt.Println("✔ Подключение client MongoDB успешно")
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("× Подключение к MongoDB не удалось:", err)
	} else {
		fmt.Println("✔ Подключение MongoDB успешно (ping) ")
	}

	db = client.Database("main")
	fmt.Println("current db name " + db.Name())
}

func UpdatePlaces(places *[]model.Place) {
	var placesDB []interface{}
	for _, v := range *places {
		placesDB = append(placesDB, v)
	}
	fsinPlacesCollection := db.Collection("fsin_places")
	insertResult, err := fsinPlacesCollection.InsertMany(context.TODO(), placesDB)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("insert result:", insertResult)
}

func Places() (places []bson.M) {
	fsinPlacesCollection := db.Collection("fsin_places")
	cursor, err := fsinPlacesCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	if err = cursor.All(context.TODO(), &places); err != nil {
		fmt.Println(err)
	}
	return places
}

func Violations() (violations []bson.M) {
	violationsCollection := db.Collection("violations")
	cursor, err := violationsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	if err = cursor.All(context.TODO(), &violations); err != nil {
		fmt.Println(err)
	}
	return violations
}

func countAllViolations() int64 {
	violationsCollection := db.Collection("violations")
	count, err := violationsCollection.EstimatedDocumentCount(context.TODO(), nil)
	if err != nil {
		fmt.Println(err)
	}
	return count
}

// получение кол-ва нарушений по типу
func CountViolations() map[string]uint32 {

	var violationsTypes = [...]string{"violations_of_medical_care", "physical_impact_from_employees", "physical_impact_from_prisoners", "psychological_impact_from_employees", "psychological_impact_from_prisoners", "corruption_from_employees", "extortions_from_employees"}
	violations := make(map[string]uint32)

	//if typeOfViolation == "" {
	//	return countAllViolations()
	//}

	violationsCollection := db.Collection("violations")
	cursor, err := violationsCollection.Find(context.TODO(), bson.M{})
	defer cursor.Close(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(context.TODO()) {
		for _, vType := range violationsTypes {
			v := cursor.Current.Lookup(vType).StringValue()
			if v != "" && v != "\t" && v != "\n" && strings.ToLower(v) != "нет" && v != "Не сталкивался с нарушениями" {
				violations[vType]++
			}
		}
	}
	if err := cursor.Err(); err != nil {
		fmt.Println(err)
	}
	return violations
}
