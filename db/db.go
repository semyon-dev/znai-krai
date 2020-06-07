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

func Places() (places []model.Place) {
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

func AddCoronaViolation(violation model.CoronaViolation) {
	coronaViolations := db.Collection("corona_violations")
	insertResult, err := coronaViolations.InsertOne(context.TODO(), violation)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("InsertedID ", insertResult.InsertedID)
}

func FindPlace(position model.Position) (place model.Place, err error) {
	fsinPlacesCollection := db.Collection("fsin_places")
	err = fsinPlacesCollection.FindOne(context.TODO(), bson.M{"position": position}).Decode(&place)
	if err != nil {
		fmt.Println("MongoDB_error: ", err)
	}
	return place, err
}

func UpdateViolation(violation model.Violation) {
	var newViolation model.Violation
	violationsCollection := db.Collection("violations")
	err := violationsCollection.FindOne(context.TODO(), bson.M{"time": violation.Time, "fsin_organization": violation.FSINOrganization}).Decode(&newViolation)
	if err != nil {
		log.Println("MongoDB err!: ", err)
	}
	//fmt.Printf("newViolation: \n %+v\n", newViolation.)

	newViolation.PlacesID = violation.PlacesID
	newViolation.Approved = violation.Approved
	newViolation.Positions = violation.Positions
	fmt.Println("violation.PlacesID: ", newViolation.PlacesID)
	fmt.Println("violation.ID: ", newViolation.ID)
	fmt.Println("violation.FSINOrganization: ", newViolation.FSINOrganization)
	update := bson.M{
		"$set": newViolation,
	}
	result, err := violationsCollection.UpdateOne(context.TODO(), bson.M{"time": violation.Time, "fsin_organization": violation.FSINOrganization}, update)
	if err != nil {
		log.Println("MongoDB error!!! -> ", err)
	}
	fmt.Printf("ModifiedCount: \n %+v\n", result.ModifiedCount)
}

func Violations() (violations []model.Violation) {
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

func CoronaViolations() (coronaViolations []model.CoronaViolation) {
	coronaViolationsCollection := db.Collection("corona_violations")
	cursor, err := coronaViolationsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	if err = cursor.All(context.TODO(), &coronaViolations); err != nil {
		fmt.Println(err)
	}
	return coronaViolations
}

func CountAllViolations() int64 {
	violationsCollection := db.Collection("violations")
	count, err := violationsCollection.EstimatedDocumentCount(context.TODO(), nil)
	if err != nil {
		fmt.Println(err)
	}
	return count
}

// получение кол-ва нарушений по типу для Аналитики
func CountViolations() map[string]map[string]uint32 {

	var violationsTypes = [...]string{
		"physical_impact_from_employees",
		"physical_impact_from_prisoners",

		"psychological_impact_from_employees",
		"psychological_impact_from_prisoners",

		"corruption_from_employees",
		"extortions_from_employees",
		"extortions_from_prisoners",

		"violations_of_medical_care",
		"visits_with_relatives",
		"communication_with_relatives",
		"communication_with_lawyer",

		"salary_of_prisoners",

		"can_prisoners_submit_complaints",
	}

	var violationsCommonTypes = [...]string{
		"physical_impact",
		"psychological_impact",
		"corruption",
		"salary_of_prisoners",

		// Есть ли у заключенных возможность направлять жалобы, ходатайства и заявления в надзирающие органы и правозащитные организации?
		"submit_complaints",

		"other",
	}

	var salaryTypes = [...]string{
		"От 0 до 100 рублей",
		"От 100 до 1 000 рублей",
		"От 1 000 до 10 000 рублей",
		"Зарплата не выплачивается",
	}

	var violations = map[string]map[string]uint32{}
	for _, v := range violationsCommonTypes {
		violations[v] = map[string]uint32{}
	}

	violationsCollection := db.Collection("violations")
	cursor, err := violationsCollection.Find(context.TODO(), bson.M{})
	if cursor == nil {
		fmt.Println("cursor is nil!")
		return nil
	}
	defer cursor.Close(context.TODO())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(context.TODO()) {
		for _, vType := range violationsTypes {
			v := cursor.Current.Lookup(vType).StringValue()
			if v != "" && v != "\t" && v != "\n" && strings.ToLower(v) != "нет" && v != "Не сталкивался с нарушениями" {
				switch {
				case vType == "physical_impact_from_employees" || vType == "physical_impact_from_prisoners":
					violations["physical_impact"][vType]++
				case vType == "psychological_impact_from_employees" || vType == "psychological_impact_from_prisoners":
					violations["psychological_impact"][vType]++
				case vType == "corruption_from_employees" || vType == "extortions_from_employees" || vType == "extortions_from_prisoners":
					violations["corruption"][vType]++
				case vType == "salary_of_prisoners":
					var exist bool
					for _, vSalary := range salaryTypes {
						if v == vSalary {
							exist = true
							break
						}
					}
					if exist {
						violations["salary_of_prisoners"][v]++
					}
				default:
					violations["other"][vType]++
				}
			}
			if vType == "can_prisoners_submit_complaints" {
				if v != "" {
					violations["submit_complaints"][v]++
				}
			}
		}
	}
	if err := cursor.Err(); err != nil {
		fmt.Println(err)
	}
	return violations
}
