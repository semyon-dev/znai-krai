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

func ShortPlaces() (shortPlaces []model.ShortPlace) {
	fsinPlacesCollection := db.Collection("fsin_places")
	cursor, err := fsinPlacesCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	if err = cursor.All(context.TODO(), &shortPlaces); err != nil {
		fmt.Println(err)
	}
	return shortPlaces
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
func CountViolations() interface{} {

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

	//var violationsCommonTypes = [...]string{
	//	"physical_impact",
	//	"psychological_impact",
	//	"corruption",
	//	"salary_of_prisoners",
	//	"communication_with_relatives",
	//	"communication_with_lawyer",
	//	"visits_with_relatives",
	//
	//	// Есть ли у заключенных возможность направлять жалобы, ходатайства и заявления в надзирающие органы и правозащитные организации?
	//	"can_prisoners_submit_complaints",
	//
	//	"other",
	//}

	var visitsWithRelatives = [...]string{
		"сокращение времени свиданий",
		"несвоевременной предоставление свиданий",
		"отказ в предоставлении свиданий",
	}

	var communicationWithOthers = [...]string{
		"отказ в телефонных звонках",
		"отказ в почтовой (телеграфной) переписке",
		"отказ в приёме передач", // отказ в приёме передач более 20 кг в колонии-поселении
		"не сталкивался с нарушениями",
		"не сразу отдают посылки",
		"цензура",                      // цензура переписки
		"нарушение конфиденциальности", // нарушение конфиденциальности свидания
		"отказ в свидании",             // отказ в свидании с заключенным
		"много нервотрепки и унижений со стороны администрации",
		"затягивание предоставления свиданий",
		"недопуск адвоката",
		"Следователь беспрепятственно может устроить допрос без адвоката",
		"ограничение времени",
	}

	var salaryTypes = [...]string{
		"От 0 до 100 рублей",
		"От 100 до 1 000 рублей",
		"От 1 000 до 10 000 рублей",
		"Зарплата не выплачивается",
	}

	//type Subcategory struct {
	//	TotalCount uint32            `json:"total_count"`
	//	Values     map[string]uint32 `json:"values"`
	//}
	//
	//type Category struct {
	//	TotalCount  uint32                 `json:"total_count"`
	//	Subcategory map[string]Subcategory `json:"subcategory"`
	//}

	type Stats struct {
		TotalCount     uint32 `json:"total_count"`
		PhysicalImpact struct {
			TotalCount                  uint32            `json:"total_count"`
			PhysicalImpactFromEmployees map[string]uint32 `json:"physical_impact_from_employees"`
			PhysicalImpactFromPrisoners map[string]uint32 `json:"physical_impact_from_prisoners"`
		} `json:"physical_impact"`
		PsychologicalImpact struct {
			TotalCount                       uint32            `json:"total_count"`
			PsychologicalImpactFromEmployees map[string]uint32 `json:"psychological_impact_from_employees"`
			PsychologicalImpactFromPrisoners map[string]uint32 `json:"psychological_impact_from_prisoners"`
		} `json:"psychological_impact"`
		Job struct {
			TotalCount        uint32            `json:"total_count"`
			LaborSlavery      uint32            `json:"labor_slavery"`
			SalaryOfPrisoners map[string]uint32 `json:"salary_of_prisoners"`
		} `json:"job"`
	}

	var stats Stats
	stats.PhysicalImpact.PhysicalImpactFromEmployees = make(map[string]uint32)
	stats.PhysicalImpact.PhysicalImpactFromPrisoners = make(map[string]uint32)
	stats.PsychologicalImpact.PsychologicalImpactFromEmployees = make(map[string]uint32)
	stats.PsychologicalImpact.PsychologicalImpactFromPrisoners = make(map[string]uint32)
	stats.Job.SalaryOfPrisoners = make(map[string]uint32)

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
				stats.TotalCount++
				switch {
				case vType == "physical_impact_from_employees":
					stats.PhysicalImpact.TotalCount++
					stats.PhysicalImpact.PhysicalImpactFromEmployees["total_count"]++
				case vType == "physical_impact_from_prisoners":
					stats.PhysicalImpact.TotalCount++
					stats.PhysicalImpact.PhysicalImpactFromPrisoners["total_count"]++
				case vType == "psychological_impact_from_employees":
					stats.PsychologicalImpact.TotalCount++
					stats.PsychologicalImpact.PsychologicalImpactFromEmployees["total_count"]++
				case vType == "psychological_impact_from_prisoners":
					stats.PsychologicalImpact.TotalCount++
					stats.PsychologicalImpact.PsychologicalImpactFromPrisoners["total_count"]++
				case vType == "corruption_from_employees" || vType == "extortions_from_employees" || vType == "extortions_from_prisoners":
					// TODO
				case vType == "communication_with_relatives" || vType == "communication_with_lawyer":
					for _, typ := range communicationWithOthers {
						if strings.Contains(strings.ToLower(v), typ) {
							// TODO: violations[vType][typ]++
						}
					}
				case vType == "visits_with_relatives":
					for _, typ := range visitsWithRelatives {
						if strings.Contains(strings.ToLower(v), typ) {
							// TODO: violations[vType][typ]++
						}
					}
				case vType == "salary_of_prisoners":
					var exist bool
					for _, vSalary := range salaryTypes {
						if v == vSalary {
							exist = true
							break
						}
					}
					if exist {
						stats.Job.SalaryOfPrisoners[v]++
					}
				default:
					// TODO: violations["other"][vType]++
				}
			}
			if vType == "can_prisoners_submit_complaints" {
				if v != "" {
					// TODO:  violations["can_prisoners_submit_complaints"][v]++
				}
			}
		}
	}
	if err := cursor.Err(); err != nil {
		fmt.Println(err)
	}
	//
	//categories["physical_impact"] = *physicalImpact
	//categories["corruption"] = *corruption
	//categories["psychological_impact"] = *psychologicalImpact

	//physicalImpact.Subcategory["physical_impact_from_prisoners"] = *physicalImpactFromPrisoners
	//physicalImpact.Subcategory["physical_impact_from_employees"] = *physicalImpactFromEmployees

	return stats
}
