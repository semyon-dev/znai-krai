package db

import (
	"context"
	"fmt"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/log"
	"github.com/semyon-dev/znai-krai/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strings"
	"time"
	"unicode"
)

var db *mongo.Database

type subcategory struct {
	Name              string            `json:"-"`
	TotalCount        uint32            `json:"total_count"`
	TotalCountAppeals uint32            `json:"total_count_appeals"`
	Values            map[string]uint32 `json:"values"`
}

type category struct {
	Name              string                 `json:"-"`
	TotalCount        uint32                 `json:"total_count"`
	TotalCountAppeals uint32                 `json:"total_count_appeals"`
	CountByYears      map[string]uint32      `json:"count_by_years"`
	Subcategories     map[string]subcategory `json:"subcategories"`
}

type Stats map[string]category

func Connect() {

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb://"+config.MongoDBLogin+":"+config.MongoDBPass+"@main-shard-00-00-h6nko.mongodb.net:27017,main-shard-00-01-h6nko.mongodb.net:27017,main-shard-00-02-h6nko.mongodb.net:27017/main?ssl=true&replicaSet=main-shard-0&authSource=admin&retryWrites=true&w=majority",
	))
	if err != nil {
		log.HandleErr(err)
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
	fmt.Println("Текущая бд: " + db.Name())
}

func Places() (places []model.Place) {
	fsinPlacesCollection := db.Collection("fsin_places")
	cursor, err := fsinPlacesCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.HandleErr(err)
	}
	if err = cursor.All(context.TODO(), &places); err != nil {
		log.HandleErr(err)
	}
	return places
}

func ShortPlaces() (shortPlaces []model.ShortPlace) {
	fsinPlacesCollection := db.Collection("fsin_places")
	cursor, err := fsinPlacesCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.HandleErr(err)
	}
	if err = cursor.All(context.TODO(), &shortPlaces); err != nil {
		log.HandleErr(err)
	}
	return shortPlaces
}

func InsertCoronaViolation(violation model.CoronaViolation) {
	coronaViolations := db.Collection("corona_violations")
	insertResult, err := coronaViolations.InsertOne(context.TODO(), violation)
	if err != nil {
		log.HandleErr(err)
	}
	fmt.Println("InsertedID ", insertResult.InsertedID)
}

func InsertReport(report model.Report) {
	reportsCollection := db.Collection("reports")
	insertResult, err := reportsCollection.InsertOne(context.TODO(), report)
	if err != nil {
		log.HandleErr(err)
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
	cursor, err := violationsCollection.Find(context.TODO(), bson.M{"approved": true})
	if err != nil {
		log.HandleErr(err)
	}
	if err = cursor.All(context.TODO(), &violations); err != nil {
		log.HandleErr(err)
	}
	return violations
}

func CoronaViolations() (coronaViolations []model.CoronaViolation) {
	coronaViolationsCollection := db.Collection("corona_violations")
	cursor, err := coronaViolationsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.HandleErr(err)
	}
	if err = cursor.All(context.TODO(), &coronaViolations); err != nil {
		log.HandleErr(err)
	}
	return coronaViolations
}

func CountAllViolations() int64 {
	violationsCollection := db.Collection("violations")
	count, err := violationsCollection.EstimatedDocumentCount(context.TODO(), nil)
	if err != nil {
		log.HandleErr(err)
	}
	return count
}

// кол-во документов (нарушений) по коронавирусу
func CountCoronaViolations() int64 {
	violationsCollection := db.Collection("corona_violations")
	count, err := violationsCollection.EstimatedDocumentCount(context.TODO(), nil)
	if err != nil {
		log.HandleErr(err)
	}
	return count
}

// Подсчет статистики нарушений для всех типов
func CountViolationsStats() (stats Stats, totalCount uint64) {

	stats = Stats{}

	violationsCollection := db.Collection("violations")
	cursor, err := violationsCollection.Find(context.TODO(), bson.M{})
	if cursor == nil {
		log.HandleErr(err)
		return stats, 0
	}
	defer func() {
		err = cursor.Close(context.TODO())
		if err != nil {
			log.HandleErr(err)
		}
	}()

	initCategory := func(categoryName string) category {
		return category{Name: categoryName, Subcategories: map[string]subcategory{}, CountByYears: map[string]uint32{}}
	}

	initSubcategory := func(subcategoryName string) subcategory {
		return subcategory{Name: subcategoryName, Values: map[string]uint32{}}
	}

	// категория physicalImpact
	var physicalImpact = initCategory("physical_impact")
	var physicalImpactFromEmployees = initSubcategory("physical_impact_from_employees")
	var physicalImpactFromPrisoners = initSubcategory("physical_impact_from_prisoners")

	// категория psychologicalImpact
	var psychologicalImpact = initCategory("psychological_impact")
	var psychologicalImpactFromPrisoners = initSubcategory("psychological_impact_from_prisoners")
	var psychologicalImpactFromEmployees = initSubcategory("psychological_impact_from_employees")

	// категория Job
	var job = initCategory("job")
	var laborSlavery = initSubcategory("labor_slavery")
	var salaryOfPrisoners = initSubcategory("salary_of_prisoners")

	// категория Corruption
	var corruption = initCategory("corruption")
	var corruptionFromEmployees = initSubcategory("corruption_from_employees")
	var extortionsFromEmployees = initSubcategory("extortions_from_employees")
	var extortionsFromPrisoners = initSubcategory("extortions_from_prisoners")

	// категория Communication
	var communication = initCategory("communication")
	var visitsWithRelatives = initSubcategory("visits_with_relatives")
	var communicationWithRelatives = initSubcategory("communication_with_relatives")
	var communicationWithLawyer = initSubcategory("communication_with_lawyer")
	var canPrisonersSubmitComplaints = initSubcategory("can_prisoners_submit_complaints")

	// категория ViolationsOfClothes
	var violationsOfClothes = initCategory("violations_of_clothes")
	var violationsOfClothesSub = initSubcategory("violations_of_clothes")

	//	категория ViolationsOfFood
	var violationsOfFood = initCategory("violations_of_food")
	var violationsOfFoodSub = initSubcategory("violations_of_food")

	//категория	ViolationsOfMedicalCare
	var violationsOfMedicalCare = initCategory("violations_of_medical_care")
	var violationsOfMedicalCareSub = initSubcategory("violations_of_medical_care")

	//категория	ViolationsStaging
	var violationsStaging = initCategory("violations_staging")
	var violationsStagingSub = initSubcategory("violations_staging")

	//	категорияReligion
	var religion = initCategory("religion")
	var violationsReligiousRitesFromEmployees = initSubcategory("violations_religious_rites_from_employees")
	var violationsReligiousRitesFromPrisoners = initSubcategory("violations_religious_rites_from_prisoners")

	// категория ViolationsWithPlacementInPunishmentCell
	var violationsWithPlacementInPunishmentCell = initCategory("violations_with_placement_in_punishment_cell")
	var violationsWithPlacementInPunishmentCellSub = initSubcategory("violations_with_placement_in_punishment_cell")

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(context.TODO()) {
		for _, vType := range model.ViolationsTypes {
			v := cursor.Current.Lookup(vType).StringValue()
			timeOfOffence := cursor.Current.Lookup("time_of_offence").StringValue()
			if v != "" && v != "\t" && v != "\n" && strings.ToLower(v) != "нет" && v != "Не сталкивался с нарушениями" {
				totalCount++
				switch vType {
				case "physical_impact_from_employees":
					countTimeOfOffence(physicalImpact.CountByYears, timeOfOffence)
					stats.countStats(&physicalImpact, &physicalImpactFromEmployees, v, model.ViolationsPhysicalImpactTypes)
				case "physical_impact_from_prisoners":
					countTimeOfOffence(physicalImpact.CountByYears, timeOfOffence)
					stats.countStats(&physicalImpact, &physicalImpactFromPrisoners, v, model.ViolationsPhysicalImpactTypes)
				case "psychological_impact_from_employees":
					countTimeOfOffence(psychologicalImpact.CountByYears, timeOfOffence)
					stats.countStats(&psychologicalImpact, &psychologicalImpactFromEmployees, v, model.ViolationsPsychologicalImpact)
				case "psychological_impact_from_prisoners":
					countTimeOfOffence(psychologicalImpact.CountByYears, timeOfOffence)
					stats.countStats(&psychologicalImpact, &psychologicalImpactFromPrisoners, v, model.ViolationsPsychologicalImpact)
				case "extortions_from_employees":
					countTimeOfOffence(corruption.CountByYears, timeOfOffence)
					stats.countStats(&corruption, &extortionsFromEmployees, v, model.ViolationsExtortionsFromEmployeesTypes)
				case "communication_with_relatives":
					countTimeOfOffence(communication.CountByYears, timeOfOffence)
					stats.countStats(&communication, &communicationWithRelatives, v, model.ViolationsCommunicationWithOthers)
				case "communication_with_lawyer":
					countTimeOfOffence(communication.CountByYears, timeOfOffence)
					stats.countStats(&communication, &communicationWithLawyer, v, model.ViolationsCommunicationWithOthers)
				case "visits_with_relatives":
					countTimeOfOffence(communication.CountByYears, timeOfOffence)
					stats.countStats(&communication, &visitsWithRelatives, v, model.ViolationsVisitsWithRelatives)
				case "violations_penalties_related_to_placement":
					countTimeOfOffence(violationsWithPlacementInPunishmentCell.CountByYears, timeOfOffence)
					stats.countStats(&violationsWithPlacementInPunishmentCell, &violationsWithPlacementInPunishmentCellSub, v, model.ViolationsWithPlacementInPunishmentCellTypes)
				case "violations_of_clothes":
					countTimeOfOffence(violationsOfClothes.CountByYears, timeOfOffence)
					stats.countStats(&violationsOfClothes, &violationsOfClothesSub, v, model.ViolationsClothes)
				case "violations_of_food":
					countTimeOfOffence(violationsOfFood.CountByYears, timeOfOffence)
					stats.countStats(&violationsOfFood, &violationsOfFoodSub, v, model.ViolationsFoodTypes)
				case "violations_of_medical_care":
					countTimeOfOffence(violationsOfMedicalCare.CountByYears, timeOfOffence)
					stats.countStats(&violationsOfMedicalCare, &violationsOfMedicalCareSub, v, model.ViolationsMedicalCareTypes)
				case "violations_religious_rites_from_employees":
					countTimeOfOffence(religion.CountByYears, timeOfOffence)
					stats.countStats(&religion, &violationsReligiousRitesFromEmployees, v, model.ViolationsReligiousViolations)
				case "violations_religious_rites_from_prisoners":
					countTimeOfOffence(religion.CountByYears, timeOfOffence)
					stats.countStats(&religion, &violationsReligiousRitesFromPrisoners, v, model.ViolationsReligiousViolations)
				case "violations_staging":
					countTimeOfOffence(violationsStaging.CountByYears, timeOfOffence)
					stats.countStats(&violationsStaging, &violationsStagingSub, v, model.ViolationsStagingViolations)
				case "salary_of_prisoners":
					countTimeOfOffence(job.CountByYears, timeOfOffence)
					stats.countStats(&job, &salaryOfPrisoners, v, model.ViolationsSalaryTypes)
				}
			}
			// если требуется ответы "Да" и "Нет" то минуем проверку на "Нет" в if перед верхним switch
			if vType == "can_prisoners_submit_complaints" && v != "" {
				countTimeOfOffence(communication.CountByYears, timeOfOffence)
				canPrisonersSubmitComplaints.Values[v]++
				canPrisonersSubmitComplaints.TotalCountAppeals++
				communication.TotalCountAppeals++
				if strings.ToLower(v) == "нет" {
					totalCount++
					communication.TotalCount++
					canPrisonersSubmitComplaints.TotalCount++
				}
				stats[communication.Name] = communication
				communication.Subcategories[canPrisonersSubmitComplaints.Name] = canPrisonersSubmitComplaints
			} else {
				if strings.ToLower(v) == "да" {
					totalCount++
				}
			}

			if vType == "corruption_from_employees" {
				countTimeOfOffence(corruption.CountByYears, timeOfOffence)
				stats.countYesNotDifficult(&corruption, &corruptionFromEmployees, v)
			} else if vType == "extortions_from_prisoners" {
				countTimeOfOffence(corruption.CountByYears, timeOfOffence)
				stats.countYesNotDifficult(&corruption, &extortionsFromPrisoners, v)
			} else if vType == "labor_slavery" {
				countTimeOfOffence(job.CountByYears, timeOfOffence)
				stats.countYesNotDifficult(&job, &laborSlavery, v)
			}
		}
	}
	if err := cursor.Err(); err != nil {
		log.HandleErr(err)
	}
	return stats, totalCount
}

// Подсчет кол-во информация по годам (map["2020"] = 23)
func countTimeOfOffence(count map[string]uint32, timeOfOffence string) {
	var data = [4]byte{}
	i := 0
	for _, symbol := range timeOfOffence {
		if unicode.IsDigit(symbol) {
			data[i] = byte(symbol)
			i++
		} else {
			if data[0] != 0 {
				count[string(data[:])]++
			}
			data = [4]byte{}
			i = 0
		}
	}
	return
}

// Подсчет статистики нарушений для аналитики
// Для заданной category и subcategory считает value пробегаясь по violationsTypes
func (stats Stats) countStats(category *category, subcategory *subcategory, value string, violationsTypes []string) {
	category.TotalCountAppeals++
	subcategory.TotalCountAppeals++
	for _, typ := range violationsTypes {
		if strings.Contains(strings.ToLower(value), typ) {
			category.TotalCount++
			subcategory.TotalCount++
			subcategory.Values[typ]++
		}
	}
	stats[category.Name] = *category
	category.Subcategories[subcategory.Name] = *subcategory
}

// Подсчет статистики нарушений для аналитики
// Вопросы где ответы "да", "нет", "затрудняюсь ответить"
func (stats Stats) countYesNotDifficult(category *category, subcategory *subcategory, value string) {
	category.TotalCountAppeals++
	subcategory.TotalCountAppeals++
	vLower := strings.ToLower(value)
	if vLower == "да" || vLower == "нет" || vLower == "затрудняюсь ответить" {
		if vLower == "да" {
			subcategory.TotalCount++
			category.TotalCount++
		}
		subcategory.Values[value]++
	}
	stats[category.Name] = *category
	category.Subcategories[subcategory.Name] = *subcategory
}
