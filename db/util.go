package db

import (
	"context"
	"fmt"
	"github.com/semyon-dev/znai-krai/model"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func UpdatePlace(place model.Place) {
	var newPlace model.Place
	fsinPlacesCollection := db.Collection("fsin_places")
	err := fsinPlacesCollection.FindOne(context.TODO(), bson.M{"_id": place.ID}).Decode(&newPlace)
	if err != nil {
		log.Fatal("MongoDB err!: ", err)
	}
	newPlace.Type = place.Type
	fmt.Println("violation.ID: ", newPlace.ID)
	fmt.Println("violation.Type: ", newPlace.Type)
	fmt.Println("---------------------------------------")
	update := bson.M{
		"$set": newPlace,
	}
	result, err := fsinPlacesCollection.UpdateOne(context.TODO(), bson.M{"_id": newPlace.ID}, update)
	if err != nil {
		log.Fatal("MongoDB error!!! -> ", err)
	}
	fmt.Printf("ModifiedCount: \n %+v\n", result.ModifiedCount)
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
