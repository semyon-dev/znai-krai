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

		"visits_with_relatives",
		"communication_with_relatives",
		"communication_with_lawyer",
		"can_prisoners_submit_complaints",

		"salary_of_prisoners",
		"labor_slavery",

		"violations_of_food",
		"violations_of_medical_care",
		"violations_of_clothes",

		"violations_staging",
		"violations_religious_rites_from_employees",
		"violations_religious_rites_from_prisoners",
		"violations_penalties_related_to_placement",
	}

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
		"следователь беспрепятственно может устроить допрос без адвоката",
		"ограничение времени",
	}

	var salaryTypes = [...]string{
		"От 0 до 100 рублей",
		"От 100 до 1 000 рублей",
		"От 1 000 до 10 000 рублей",
		"Зарплата не выплачивается",
	}

	var ViolationsClothes = [...]string{
		"отсутствие (несвоевременная выдача) зимней одежды и обуви",
		"отсутствие одежды и обуви по размеру",
	}

	var ViolationsFood = [...]string{
		"маленькие порции",
		"отсутствие мясных продуктов",
		"испорченные продукты",
		"однообразное питание",
		"однообразное меню",
		"кислые овощи",
		"отсутствие молочных продуктов, грязные каши",
		"грязные каши",
	}

	var ViolationsMedicalCare = [...]string{
		"отказ в оказании медицинской помощи",
		"медицинская помощь была недостаточная или оказана с задержкой",
		"отказ в стоматологической помощи",
		"медикаменты отсутствовали или были испорчены",
	}

	var PsychologicalImpact = [...]string{
		"угроза применения пыток",
		"угроза применения административного взыскания",
		"унижение",
		"применение коллективного взыскания к группе заключенных",
		"угроза моим детям",
		"угроза жизни", // Угроза жизни осуждённому
		"применение силы без причины",
		"угроза закрыть в ЕПКТ",
	}

	var religiousViolations = [...]string{
		"отказ в посещении храма",
		"запрет ночной молитвы",
		"запрет на хранение (передачу) религиозной литературы",
		"предметов культа",
		"оскорбления",
		"притеснения по религиозным мотивам",
		"в браке не давали молиться активисты",
	}

	var stagingViolations = [...]string{
		"переполненность сборочной камеры",
		"отсутствие вентиляции (отопления)",
		"совместное нахождение с инфекционными больными",
		"неоказание медицинской помощи",
		"отсутствие (недостаток) питания",
		"перевозка заключенных в «стаканах»",
	}

	var physicalImpact = [...]string{
		"избиение",
		"применение наручников",
		"избиение",
		"унижение",
		"унизительные процедуры",
		"вымогательство",
		"применение электрического тока",
		"подвешивание",
		"привязывание",
		"нахождение в неудобной позе",
		"удушение",
		"провокация",
		"воруют вещи",
		"перекидывают в разные камеры",
		"лишение еды и питья",
		"поливание водой",
		"оставляли зимой в бетонной комнате с открытыми окнами и раздетым до трусов",
		"тяжелая физическая работа без ограничения продолжительности",
	}

	type Stats struct {
		TotalCount     uint32 `json:"total_count"`
		PhysicalImpact struct {
			TotalCount                  uint32            `json:"total_count"`
			TotalCountAppeals           uint32            `json:"total_count_appeals"`
			PhysicalImpactFromEmployees map[string]uint32 `json:"physical_impact_from_employees"`
			PhysicalImpactFromPrisoners map[string]uint32 `json:"physical_impact_from_prisoners"`
		} `json:"physical_impact"`
		PsychologicalImpact struct {
			TotalCount                       uint32            `json:"total_count"`
			TotalCountAppeals                uint32            `json:"total_count_appeals"`
			PsychologicalImpactFromEmployees map[string]uint32 `json:"psychological_impact_from_employees"`
			PsychologicalImpactFromPrisoners map[string]uint32 `json:"psychological_impact_from_prisoners"`
		} `json:"psychological_impact"`
		Job struct {
			TotalCount        uint32            `json:"total_count"`
			LaborSlavery      uint32            `json:"labor_slavery"`
			SalaryOfPrisoners map[string]uint32 `json:"salary_of_prisoners"`
		} `json:"job"`
		Corruption struct {
			TotalCount              uint32            `json:"total_count"`
			CorruptionFromEmployees map[string]uint32 `json:"corruption_from_employees"`
			ExtortionsFromEmployees map[string]uint32 `json:"extortions_from_employees"`
			ExtortionsFromPrisoners map[string]uint32 `json:"extortions_from_prisoners"`
		} `json:"corruption"`
		Communication struct {
			TotalCount                   uint32            `json:"total_count"`
			VisitsWithRelatives          map[string]uint32 `json:"visits_with_relatives"`
			CommunicationWithRelatives   map[string]uint32 `json:"communication_with_relatives"`
			CommunicationWithLawyer      map[string]uint32 `json:"communication_with_lawyer"`
			CanPrisonersSubmitComplaints map[string]uint32 `json:"can_prisoners_submit_complaints"`
		} `json:"communication"`
		ViolationsOfClothes struct {
			TotalCount          uint32            `json:"total_count"`
			ViolationsOfClothes map[string]uint32 `json:"violations_of_clothes"`
		} `json:"violations_of_clothes"`
		ViolationsOfFood struct {
			TotalCount       uint32            `json:"total_count"`
			ViolationsOfFood map[string]uint32 `json:"violations_of_food"`
		} `json:"violations_of_food"`
		ViolationsOfMedicalCare struct {
			TotalCount              uint32            `json:"total_count"`
			ViolationsOfMedicalCare map[string]uint32 `json:"violations_of_medical_care"`
		} `json:"violations_of_medical_care"`
		ViolationsStaging struct {
			TotalCount        uint32            `json:"total_count"`
			ViolationsStaging map[string]uint32 `json:"violations_staging"`
		} `json:"violations_staging"`
		Religion struct {
			TotalCount                            uint32            `json:"total_count"`
			ViolationsReligiousRitesFromEmployees map[string]uint32 `json:"violations_religious_rites_from_employees"`
			ViolationsReligiousRitesFromPrisoners map[string]uint32 `json:"violations_religious_rites_from_prisoners"`
		} `json:"religion"`
	}

	var stats Stats

	stats.PhysicalImpact.PhysicalImpactFromEmployees = make(map[string]uint32)
	stats.PhysicalImpact.PhysicalImpactFromPrisoners = make(map[string]uint32)

	stats.PsychologicalImpact.PsychologicalImpactFromEmployees = make(map[string]uint32)
	stats.PsychologicalImpact.PsychologicalImpactFromPrisoners = make(map[string]uint32)

	stats.Job.SalaryOfPrisoners = make(map[string]uint32)

	stats.Corruption.CorruptionFromEmployees = make(map[string]uint32)
	stats.Corruption.ExtortionsFromEmployees = make(map[string]uint32)
	stats.Corruption.ExtortionsFromPrisoners = make(map[string]uint32)

	stats.Communication.VisitsWithRelatives = make(map[string]uint32)
	stats.Communication.CommunicationWithLawyer = make(map[string]uint32)
	stats.Communication.CommunicationWithRelatives = make(map[string]uint32)
	stats.Communication.CanPrisonersSubmitComplaints = make(map[string]uint32)

	stats.ViolationsOfClothes.ViolationsOfClothes = make(map[string]uint32)
	stats.ViolationsOfFood.ViolationsOfFood = make(map[string]uint32)
	stats.ViolationsOfMedicalCare.ViolationsOfMedicalCare = make(map[string]uint32)

	stats.Religion.ViolationsReligiousRitesFromEmployees = make(map[string]uint32)
	stats.Religion.ViolationsReligiousRitesFromPrisoners = make(map[string]uint32)

	stats.ViolationsStaging.ViolationsStaging = make(map[string]uint32)

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
				switch vType {
				case "physical_impact_from_employees":
					stats.PhysicalImpact.TotalCountAppeals++
					stats.PhysicalImpact.PhysicalImpactFromEmployees["total_count"]++
					for _, typ := range physicalImpact {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.PhysicalImpact.TotalCount++
							stats.PhysicalImpact.PhysicalImpactFromEmployees[typ]++
						}
					}
				case "physical_impact_from_prisoners":
					stats.PhysicalImpact.TotalCountAppeals++
					stats.PhysicalImpact.PhysicalImpactFromPrisoners["total_count"]++
					for _, typ := range physicalImpact {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.PhysicalImpact.TotalCount++
							stats.PhysicalImpact.PhysicalImpactFromPrisoners[typ]++
						}
					}
				case "psychological_impact_from_employees":
					stats.PsychologicalImpact.TotalCountAppeals++
					stats.PsychologicalImpact.PsychologicalImpactFromEmployees["total_count"]++
					for _, typ := range PsychologicalImpact {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.PsychologicalImpact.TotalCount++
							stats.PsychologicalImpact.PsychologicalImpactFromEmployees[typ]++
						}
					}
				case "psychological_impact_from_prisoners":
					stats.PsychologicalImpact.TotalCountAppeals++
					stats.PsychologicalImpact.PsychologicalImpactFromPrisoners["total_count"]++
					for _, typ := range PsychologicalImpact {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.PsychologicalImpact.TotalCount++
							stats.PsychologicalImpact.PsychologicalImpactFromPrisoners[typ]++
						}
					}
				case "corruption_from_employees":
					stats.Corruption.TotalCount++
					stats.Corruption.CorruptionFromEmployees["total_count"]++
				case "extortions_from_prisoners":
					stats.Corruption.TotalCount++
					stats.Corruption.ExtortionsFromPrisoners["total_count"]++
				case "extortions_from_employees":
					stats.Corruption.TotalCount++
					stats.Corruption.ExtortionsFromEmployees["total_count"]++
				case "communication_with_relatives":
					for _, typ := range communicationWithOthers {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.Communication.TotalCount++
							stats.Communication.CommunicationWithRelatives["total_count"]++
							stats.Communication.CommunicationWithRelatives[typ]++
						}
					}
				case "communication_with_lawyer":
					for _, typ := range communicationWithOthers {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.Communication.TotalCount++
							stats.Communication.CommunicationWithLawyer["total_count"]++
							stats.Communication.CommunicationWithLawyer[typ]++
						}
					}
				case "visits_with_relatives":
					for _, typ := range visitsWithRelatives {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.Communication.TotalCount++
							stats.Communication.VisitsWithRelatives["total_count"]++
							stats.Communication.VisitsWithRelatives[typ]++
						}
					}
				case "salary_of_prisoners":
					var exist bool
					for _, vSalary := range salaryTypes {
						if v == vSalary {
							exist = true
							break
						}
					}
					if exist {
						stats.Job.TotalCount++
						stats.Job.SalaryOfPrisoners["total_count"]++
						stats.Job.SalaryOfPrisoners[v]++
					}
				case "violations_of_clothes":
					stats.ViolationsOfClothes.ViolationsOfClothes["total_count_appeals"]++
					for _, typ := range ViolationsClothes {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.ViolationsOfClothes.TotalCount++
							stats.ViolationsOfClothes.ViolationsOfClothes[typ]++
						}
					}
				case "violations_of_food":
					stats.ViolationsOfFood.ViolationsOfFood["total_count_appeals"]++
					for _, typ := range ViolationsFood {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.ViolationsOfFood.TotalCount++
							stats.ViolationsOfFood.ViolationsOfFood[typ]++
						}
					}
				case "violations_of_medical_care":
					stats.ViolationsOfMedicalCare.ViolationsOfMedicalCare["total_count_appeals"]++
					for _, typ := range ViolationsMedicalCare {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.ViolationsOfMedicalCare.TotalCount++
							stats.ViolationsOfMedicalCare.ViolationsOfMedicalCare[typ]++
						}
					}
				case "violations_religious_rites_from_employees":
					stats.Religion.TotalCount++
					for _, typ := range religiousViolations {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.Religion.ViolationsReligiousRitesFromEmployees["total_count"]++
							stats.Religion.ViolationsReligiousRitesFromEmployees[typ]++
						}
					}
				case "violations_religious_rites_from_prisoners":
					stats.Religion.TotalCount++
					for _, typ := range religiousViolations {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.Religion.ViolationsReligiousRitesFromPrisoners["total_count"]++
							stats.Religion.ViolationsReligiousRitesFromPrisoners[typ]++
						}
					}
				case "violations_staging":
					for _, typ := range stagingViolations {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.ViolationsStaging.TotalCount++
							stats.ViolationsStaging.ViolationsStaging[typ]++
						}
					}
				}
			}
			if vType == "can_prisoners_submit_complaints" && v != "" {
				if strings.ToLower(v) == "нет" {
					stats.TotalCount++
					stats.Communication.TotalCount++
				}
				stats.Communication.CanPrisonersSubmitComplaints[v]++
				stats.Communication.CanPrisonersSubmitComplaints["total_count"]++
			}
		}
	}
	if err := cursor.Err(); err != nil {
		fmt.Println(err)
	}
	return stats
}
