package db

import (
	"context"
	"fmt"
	"github.com/semyon-dev/znai-krai/log"
	"github.com/semyon-dev/znai-krai/model"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func UpdateViolation(violation model.Violation) {
	var newViolation model.Violation
	violationsCollection := db.Collection("violations")
	err := violationsCollection.FindOne(context.TODO(), bson.M{"time": violation.Time, "fsin_organization": violation.FSINOrganization}).Decode(&newViolation)
	if err != nil {
		log.HandleErr(err)
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
		log.HandleErr(err)
	}
	fmt.Printf("ModifiedCount: \n %+v\n", result.ModifiedCount)
}

func UpdatePlace(place model.Place) {
	var newPlace model.Place
	fsinPlacesCollection := db.Collection("fsin_places")
	err := fsinPlacesCollection.FindOne(context.TODO(), bson.M{"name": place.Name, "address": place.Address, "location": place.Location}).Decode(&newPlace)
	if err != nil {
		log.HandleErr(err)
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
		log.HandleErr(err)
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
		log.HandleErr(err)
	}
	fmt.Println("insert result:", insertResult)
}

// получение кол-ва нарушений по типу для Аналитики
// Deprecated:
func CountViolationsOld() interface{} {

	type PhysicalImpact struct {
		CountByYears                map[string]uint32 `json:"count_by_years"`
		TotalCountAppeals           uint32            `json:"total_count_appeals"`
		TotalCount                  uint32            `json:"total_count"`
		PhysicalImpactFromEmployees map[string]uint32 `json:"physical_impact_from_employees"`
		PhysicalImpactFromPrisoners map[string]uint32 `json:"physical_impact_from_prisoners"`
	}

	type Stats struct {
		TotalCount          uint32 `json:"total_count"`
		PhysicalImpact      `json:"physical_impact"`
		PsychologicalImpact struct {
			CountByYears                     map[string]uint32 `json:"count_by_years"`
			TotalCountAppeals                uint32            `json:"total_count_appeals"`
			TotalCount                       uint32            `json:"total_count"`
			PsychologicalImpactFromEmployees map[string]uint32 `json:"psychological_impact_from_employees"`
			PsychologicalImpactFromPrisoners map[string]uint32 `json:"psychological_impact_from_prisoners"`
		} `json:"psychological_impact"`
		Job struct {
			CountByYears      map[string]uint32 `json:"count_by_years"`
			TotalCountAppeals uint32            `json:"total_count_appeals"`
			TotalCount        uint32            `json:"total_count"`
			LaborSlavery      map[string]uint32 `json:"labor_slavery"`
			SalaryOfPrisoners map[string]uint32 `json:"salary_of_prisoners"`
		} `json:"job"`
		Corruption struct {
			CountByYears            map[string]uint32 `json:"count_by_years"`
			TotalCountAppeals       uint32            `json:"total_count_appeals"`
			TotalCount              uint32            `json:"total_count"`
			CorruptionFromEmployees map[string]uint32 `json:"corruption_from_employees"`
			ExtortionsFromEmployees map[string]uint32 `json:"extortions_from_employees"`
			ExtortionsFromPrisoners map[string]uint32 `json:"extortions_from_prisoners"`
		} `json:"corruption"`
		Communication struct {
			CountByYears                 map[string]uint32 `json:"count_by_years"`
			TotalCountAppeals            uint32            `json:"total_count_appeals"`
			TotalCount                   uint32            `json:"total_count"`
			VisitsWithRelatives          map[string]uint32 `json:"visits_with_relatives"`
			CommunicationWithRelatives   map[string]uint32 `json:"communication_with_relatives"`
			CommunicationWithLawyer      map[string]uint32 `json:"communication_with_lawyer"`
			CanPrisonersSubmitComplaints map[string]uint32 `json:"can_prisoners_submit_complaints"`
		} `json:"communication"`
		ViolationsOfClothes struct {
			CountByYears        map[string]uint32 `json:"count_by_years"`
			TotalCountAppeals   uint32            `json:"total_count_appeals"`
			TotalCount          uint32            `json:"total_count"`
			ViolationsOfClothes map[string]uint32 `json:"violations_of_clothes"`
		} `json:"violations_of_clothes"`
		ViolationsOfFood struct {
			CountByYears      map[string]uint32 `json:"count_by_years"`
			TotalCountAppeals uint32            `json:"total_count_appeals"`
			TotalCount        uint32            `json:"total_count"`
			ViolationsOfFood  map[string]uint32 `json:"violations_of_food"`
		} `json:"violations_of_food"`
		ViolationsOfMedicalCare struct {
			CountByYears            map[string]uint32 `json:"count_by_years"`
			TotalCountAppeals       uint32            `json:"total_count_appeals"`
			TotalCount              uint32            `json:"total_count"`
			ViolationsOfMedicalCare map[string]uint32 `json:"violations_of_medical_care"`
		} `json:"violations_of_medical_care"`
		ViolationsStaging struct {
			CountByYears      map[string]uint32 `json:"count_by_years"`
			TotalCountAppeals uint32            `json:"total_count_appeals"`
			TotalCount        uint32            `json:"total_count"`
			ViolationsStaging map[string]uint32 `json:"violations_staging"`
		} `json:"violations_staging"`
		Religion struct {
			CountByYears                          map[string]uint32 `json:"count_by_years"`
			TotalCountAppeals                     uint32            `json:"total_count_appeals"`
			TotalCount                            uint32            `json:"total_count"`
			ViolationsReligiousRitesFromEmployees map[string]uint32 `json:"violations_religious_rites_from_employees"`
			ViolationsReligiousRitesFromPrisoners map[string]uint32 `json:"violations_religious_rites_from_prisoners"`
		} `json:"religion"`
		ViolationsWithPlacementInPunishmentCell struct {
			CountByYears                            map[string]uint32 `json:"count_by_years"`
			TotalCountAppeals                       uint32            `json:"total_count_appeals"`
			TotalCount                              uint32            `json:"total_count"`
			ViolationsWithPlacementInPunishmentCell map[string]uint32 `json:"violations_with_placement_in_punishment_cell"`
		} `json:"violations_with_placement_in_punishment_cell"`
	}

	var stats Stats

	stats.PhysicalImpact.PhysicalImpactFromEmployees = make(map[string]uint32)
	stats.PhysicalImpact.PhysicalImpactFromPrisoners = make(map[string]uint32)
	stats.PhysicalImpact.CountByYears = make(map[string]uint32)

	stats.PsychologicalImpact.PsychologicalImpactFromEmployees = make(map[string]uint32)
	stats.PsychologicalImpact.PsychologicalImpactFromPrisoners = make(map[string]uint32)
	stats.PsychologicalImpact.CountByYears = make(map[string]uint32)

	stats.Job.SalaryOfPrisoners = make(map[string]uint32)
	stats.Job.LaborSlavery = make(map[string]uint32)
	stats.Job.CountByYears = make(map[string]uint32)

	stats.Corruption.CorruptionFromEmployees = make(map[string]uint32)
	stats.Corruption.ExtortionsFromEmployees = make(map[string]uint32)
	stats.Corruption.ExtortionsFromPrisoners = make(map[string]uint32)
	stats.Corruption.CountByYears = make(map[string]uint32)

	stats.Communication.VisitsWithRelatives = make(map[string]uint32)
	stats.Communication.CommunicationWithLawyer = make(map[string]uint32)
	stats.Communication.CommunicationWithRelatives = make(map[string]uint32)
	stats.Communication.CanPrisonersSubmitComplaints = make(map[string]uint32)
	stats.Communication.CountByYears = make(map[string]uint32)

	stats.ViolationsOfClothes.ViolationsOfClothes = make(map[string]uint32)
	stats.ViolationsOfClothes.CountByYears = make(map[string]uint32)

	stats.ViolationsOfFood.ViolationsOfFood = make(map[string]uint32)
	stats.ViolationsOfFood.CountByYears = make(map[string]uint32)

	stats.ViolationsOfMedicalCare.ViolationsOfMedicalCare = make(map[string]uint32)
	stats.ViolationsOfMedicalCare.CountByYears = make(map[string]uint32)

	stats.Religion.ViolationsReligiousRitesFromEmployees = make(map[string]uint32)
	stats.Religion.ViolationsReligiousRitesFromPrisoners = make(map[string]uint32)
	stats.Religion.CountByYears = make(map[string]uint32)

	stats.ViolationsStaging.ViolationsStaging = make(map[string]uint32)
	stats.ViolationsStaging.CountByYears = make(map[string]uint32)
	stats.ViolationsWithPlacementInPunishmentCell.ViolationsWithPlacementInPunishmentCell = make(map[string]uint32)
	stats.ViolationsWithPlacementInPunishmentCell.CountByYears = make(map[string]uint32)

	violationsCollection := db.Collection("violations")
	cursor, err := violationsCollection.Find(context.TODO(), bson.M{})
	if cursor == nil {
		log.HandleErr(err)
		return nil
	}
	defer cursor.Close(context.TODO())
	if err != nil {
		log.HandleErr(err)
		return nil
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(context.TODO()) {
		for _, vType := range model.ViolationsTypes {
			v := cursor.Current.Lookup(vType).StringValue()
			timeOfOffence := cursor.Current.Lookup("time_of_offence").StringValue()
			if v != "" && v != "\t" && v != "\n" && strings.ToLower(v) != "нет" && v != "Не сталкивался с нарушениями" {
				stats.TotalCount++
				switch vType {
				case "physical_impact_from_employees":
					countTimeOfOffence(stats.PhysicalImpact.CountByYears, timeOfOffence)
					stats.PhysicalImpact.TotalCountAppeals++
					stats.PhysicalImpact.PhysicalImpactFromEmployees["total_count_appeals"]++
					for _, typ := range model.ViolationsPhysicalImpactTypes {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.PhysicalImpact.TotalCount++
							stats.PhysicalImpact.PhysicalImpactFromEmployees["total_count"]++
							stats.PhysicalImpact.PhysicalImpactFromEmployees[typ]++
						}
					}
				case "physical_impact_from_prisoners":
					countTimeOfOffence(stats.PhysicalImpact.CountByYears, timeOfOffence)
					stats.PhysicalImpact.TotalCountAppeals++
					stats.PhysicalImpact.PhysicalImpactFromPrisoners["total_count_appeals"]++
					for _, typ := range model.ViolationsPhysicalImpactTypes {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.PhysicalImpact.TotalCount++
							stats.PhysicalImpact.PhysicalImpactFromPrisoners["total_count"]++
							stats.PhysicalImpact.PhysicalImpactFromPrisoners[typ]++
						}
					}
				case "psychological_impact_from_employees":
					countTimeOfOffence(stats.PsychologicalImpact.CountByYears, timeOfOffence)
					stats.PsychologicalImpact.TotalCountAppeals++
					stats.PsychologicalImpact.PsychologicalImpactFromEmployees["total_count_appeals"]++
					for _, typ := range model.ViolationsPsychologicalImpact {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.PsychologicalImpact.TotalCount++
							stats.PsychologicalImpact.PsychologicalImpactFromEmployees["total_count"]++
							stats.PsychologicalImpact.PsychologicalImpactFromEmployees[typ]++
						}
					}
				case "psychological_impact_from_prisoners":
					countTimeOfOffence(stats.PsychologicalImpact.CountByYears, timeOfOffence)
					stats.PsychologicalImpact.TotalCountAppeals++
					stats.PsychologicalImpact.PsychologicalImpactFromPrisoners["total_count_appeals"]++
					for _, typ := range model.ViolationsPsychologicalImpact {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.PsychologicalImpact.TotalCount++
							stats.PsychologicalImpact.PsychologicalImpactFromPrisoners["total_count"]++
							stats.PsychologicalImpact.PsychologicalImpactFromPrisoners[typ]++
						}
					}
				case "extortions_from_employees":
					countTimeOfOffence(stats.Corruption.CountByYears, timeOfOffence)
					stats.Corruption.TotalCountAppeals++
					stats.Corruption.ExtortionsFromEmployees["total_count_appeals"]++
					for _, typ := range model.ViolationsExtortionsFromEmployeesTypes {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.Corruption.TotalCount++
							stats.Corruption.ExtortionsFromEmployees["total_count"]++
							stats.Corruption.ExtortionsFromEmployees[typ]++
						}
					}
				case "communication_with_relatives":
					countTimeOfOffence(stats.Communication.CountByYears, timeOfOffence)
					stats.Communication.TotalCountAppeals++
					stats.Communication.CommunicationWithRelatives["total_count_appeals"]++
					for _, typ := range model.ViolationsCommunicationWithOthers {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.Communication.TotalCount++
							stats.Communication.CommunicationWithRelatives["total_count"]++
							stats.Communication.CommunicationWithRelatives[typ]++
						}
					}
				case "communication_with_lawyer":
					countTimeOfOffence(stats.Communication.CountByYears, timeOfOffence)
					stats.Communication.TotalCountAppeals++
					stats.Communication.CommunicationWithLawyer["total_count_appeals"]++
					for _, typ := range model.ViolationsCommunicationWithOthers {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.Communication.TotalCount++
							stats.Communication.CommunicationWithLawyer["total_count"]++
							stats.Communication.CommunicationWithLawyer[typ]++
						}
					}
				case "visits_with_relatives":
					countTimeOfOffence(stats.Communication.CountByYears, timeOfOffence)
					stats.Communication.TotalCountAppeals++
					stats.Communication.VisitsWithRelatives["total_count_appeals"]++
					for _, typ := range model.ViolationsVisitsWithRelatives {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.Communication.TotalCount++
							stats.Communication.VisitsWithRelatives["total_count"]++
							stats.Communication.VisitsWithRelatives[typ]++
						}
					}
				case "violations_penalties_related_to_placement":
					countTimeOfOffence(stats.ViolationsWithPlacementInPunishmentCell.CountByYears, timeOfOffence)
					stats.ViolationsWithPlacementInPunishmentCell.TotalCountAppeals++
					stats.ViolationsWithPlacementInPunishmentCell.ViolationsWithPlacementInPunishmentCell["total_count_appeals"]++
					for _, typ := range model.ViolationsWithPlacementInPunishmentCellTypes {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.ViolationsWithPlacementInPunishmentCell.TotalCount++
							stats.ViolationsWithPlacementInPunishmentCell.ViolationsWithPlacementInPunishmentCell["total_count"]++
							stats.ViolationsWithPlacementInPunishmentCell.ViolationsWithPlacementInPunishmentCell[typ]++
						}
					}
				case "salary_of_prisoners":
					countTimeOfOffence(stats.Job.CountByYears, timeOfOffence)
					var exist bool
					for _, vSalary := range model.ViolationsSalaryTypes {
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
					countTimeOfOffence(stats.ViolationsOfClothes.CountByYears, timeOfOffence)
					stats.ViolationsOfClothes.TotalCountAppeals++
					stats.ViolationsOfClothes.ViolationsOfClothes["total_count_appeals"]++
					for _, typ := range model.ViolationsClothes {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.ViolationsOfClothes.TotalCount++
							stats.ViolationsOfClothes.ViolationsOfClothes["total_count"]++
							stats.ViolationsOfClothes.ViolationsOfClothes[typ]++
						}
					}
				case "violations_of_food":
					countTimeOfOffence(stats.ViolationsOfFood.CountByYears, timeOfOffence)
					stats.ViolationsOfFood.TotalCountAppeals++
					stats.ViolationsOfFood.ViolationsOfFood["total_count_appeals"]++
					for _, typ := range model.ViolationsFoodTypes {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.ViolationsOfFood.TotalCount++
							stats.ViolationsOfFood.ViolationsOfFood["total_count"]++
							stats.ViolationsOfFood.ViolationsOfFood[typ]++
						}
					}
				case "violations_of_medical_care":
					countTimeOfOffence(stats.ViolationsOfMedicalCare.CountByYears, timeOfOffence)
					stats.ViolationsOfMedicalCare.TotalCountAppeals++
					stats.ViolationsOfMedicalCare.ViolationsOfMedicalCare["total_count_appeals"]++
					for _, typ := range model.ViolationsMedicalCareTypes {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.ViolationsOfMedicalCare.TotalCount++
							stats.ViolationsOfMedicalCare.ViolationsOfMedicalCare["total_count"]++
							stats.ViolationsOfMedicalCare.ViolationsOfMedicalCare[typ]++
						}
					}
				case "violations_religious_rites_from_employees":
					countTimeOfOffence(stats.Religion.CountByYears, timeOfOffence)
					stats.Religion.TotalCountAppeals++
					stats.Religion.ViolationsReligiousRitesFromEmployees["total_count_appeals"]++
					for _, typ := range model.ViolationsReligiousViolations {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.Religion.TotalCount++
							stats.Religion.ViolationsReligiousRitesFromEmployees["total_count"]++
							stats.Religion.ViolationsReligiousRitesFromEmployees[typ]++
						}
					}
				case "violations_religious_rites_from_prisoners":
					countTimeOfOffence(stats.Religion.CountByYears, timeOfOffence)
					stats.Religion.TotalCountAppeals++
					stats.Religion.ViolationsReligiousRitesFromPrisoners["total_count_appeals"]++
					for _, typ := range model.ViolationsReligiousViolations {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.Religion.TotalCount++
							stats.Religion.ViolationsReligiousRitesFromPrisoners["total_count"]++
							stats.Religion.ViolationsReligiousRitesFromPrisoners[typ]++
						}
					}
				case "violations_staging":
					countTimeOfOffence(stats.ViolationsStaging.CountByYears, timeOfOffence)
					stats.ViolationsStaging.TotalCountAppeals++
					stats.ViolationsStaging.ViolationsStaging["total_count_appeals"]++
					for _, typ := range model.ViolationsStagingViolations {
						if strings.Contains(strings.ToLower(v), typ) {
							stats.ViolationsStaging.TotalCount++
							stats.ViolationsStaging.ViolationsStaging["total_count"]++
							stats.ViolationsStaging.ViolationsStaging[typ]++
						}
					}
				}
			}
			// если требуется ответы "Да" и "Нет" то минуем проверку на "Нет" в if перед верхним switch
			if vType == "can_prisoners_submit_complaints" && v != "" {
				countTimeOfOffence(stats.Communication.CountByYears, timeOfOffence)
				stats.Communication.CanPrisonersSubmitComplaints[v]++
				stats.Communication.CanPrisonersSubmitComplaints["total_count_appeals"]++
				stats.Communication.TotalCountAppeals++
				if strings.ToLower(v) == "нет" {
					stats.Communication.TotalCount++
					stats.Communication.CanPrisonersSubmitComplaints["total_count"]++
				}
			} else if vType == "corruption_from_employees" {
				countTimeOfOffence(stats.Corruption.CountByYears, timeOfOffence)
				stats.Corruption.TotalCountAppeals++
				stats.Corruption.CorruptionFromEmployees["total_count_appeals"]++
				vLower := strings.ToLower(v)
				if vLower == "да" || vLower == "нет" || vLower == "затрудняюсь ответить" {
					if vLower == "да" {
						stats.Corruption.CorruptionFromEmployees["total_count"]++
						stats.Corruption.TotalCount++
					}
					stats.Corruption.CorruptionFromEmployees[v]++
				}
			} else if vType == "extortions_from_prisoners" {
				countTimeOfOffence(stats.Corruption.CountByYears, timeOfOffence)
				stats.Corruption.TotalCountAppeals++
				stats.Corruption.ExtortionsFromPrisoners["total_count_appeals"]++
				vLower := strings.ToLower(v)
				if vLower == "да" || vLower == "нет" || vLower == "затрудняюсь ответить" {
					if vLower == "да" {
						stats.Corruption.TotalCount++
						stats.Corruption.ExtortionsFromPrisoners["total_count"]++
					}
					stats.Corruption.ExtortionsFromPrisoners[v]++
				}
			} else if vType == "labor_slavery" {
				countTimeOfOffence(stats.Job.CountByYears, timeOfOffence)
				stats.Job.TotalCountAppeals++
				stats.Job.LaborSlavery["total_count_appeals"]++
				vLower := strings.ToLower(v)
				if vLower == "да" || vLower == "нет" || vLower == "затрудняюсь ответить" {
					if vLower == "да" {
						stats.Job.LaborSlavery["total_count"]++
						stats.Job.TotalCount++
					}
					stats.Job.LaborSlavery[v]++
				}
			}
			if v == "Да" {
				stats.TotalCount++
			}
		}
	}
	if err := cursor.Err(); err != nil {
		log.HandleErr(err)
	}
	return stats
}
