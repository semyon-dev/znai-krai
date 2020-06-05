package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// структура для нарушений
type Form struct {
	// отметка времени
	Time string `json:"time" bson:"time"`
	// Какой Ваш статус?
	Status string `json:"status" bson:"status" binding:"required"`
	// В каком регионе находится учреждение ФСИН о котором Вы рассказали?
	Region string `json:"region" bson:"region" binding:"required"`
	// О каком учреждении ФСИН Вы рассказали?
	FSINOrganization string `json:"fsin_organization" binding:"required" bson:"fsin_organization"`
	// Укажите когда произошли нарушения о которых Вы рассказали?
	TimeOfOffence string `json:"time_of_offence" bson:"time_of_offence"`

	/// Category: физическое насилие
	// С какими фактами применения физического воздействия со стороны сотрудников ФСИН Вам приходилось сталкиваться?
	PhysicalImpactFromEmployees string `json:"physical_impact_from_employees" bson:"physical_impact_from_employees"`
	// С какими фактами применения физического воздействия со стороны заключенных Вам приходилось сталкиваться?
	PhysicalImpactFromPrisoners string `json:"physical_impact_from_prisoners" bson:"physical_impact_from_prisoners"`

	/// Category: психологическогое воздействие
	// С какими фактами психологического воздействия со стороны сотрудников ФСИН Вам приходилось сталкиваться?
	PsychologicalImpactFromEmployees string `json:"psychological_impact_from_employees" bson:"psychological_impact_from_employees"`
	// С какими фактами психологического воздействия со стороны заключенных Вам приходилось сталкиваться?
	PsychologicalImpactFromPrisoners string `json:"psychological_impact_from_prisoners" bson:"psychological_impact_from_prisoners"`

	/// Category: коррупция и вымогательства
	// В каких случаях Вы сталкивались с фактами вымогательства со стороны сотрудников ФСИН?
	ExtortionsFromEmployees string `json:"extortions_from_employees" bson:"extortions_from_employees"`
	// Приходилось ли Вам сталкиваться с иными случаями коррупции сотрудников ФСИН?
	CorruptionFromEmployees string `json:"corruption_from_employees" bson:"corruption_from_employees"`
	// Приходилось ли Вам сталкиваться с фактами вымогательства со стороны заключенных?
	ExtortionsFromPrisoners string `json:"extortions_from_prisoners" bson:"extortions_from_prisoners"`

	// Какие нарушения, связанные с оказанием медицинской помощи, Вы можете отметить?
	ViolationsOfMedicalCare string `json:"violations_of_medical_care" bson:"violations_of_medical_care"`

	// Какие нарушения, связанные с питанием, Вы можете отметить?
	ViolationsOfFood string `json:"violations_of_food" bson:"violations_of_food"`

	// Какие нарушения, связанные с одеждой, Вы можете отметить?
	ViolationsOfClothes string `json:"violations_of_clothes" bson:"violations_of_clothes"`

	// Известны ли Вам случаи трудового рабства?
	LaborSlavery string `json:"labor_slavery" bson:"labor_slavery"`

	// Укажите, в каком диапазоне находится месячная зарплата заключенных?
	SalaryOfPrisoners string `json:"salary_of_prisoners" bson:"salary_of_prisoners"`

	// Какие нарушения, связанные с предоставлением свиданий с Родственниками, Вам известны?
	VisitsWithRelatives string `json:"visits_with_relatives" bson:"visits_with_relatives"`

	// Какие нарушения, связанные с иными формами общения с Родственниками, Вам известны?
	CommunicationWithRelatives string `json:"communication_with_relatives" bson:"communication_with_relatives"`

	// Какие нарушения, связанные с общением с адвокатом (иным лицом, имеющим право на оказание юридической помощи), Вам известны?
	CommunicationWithLawyer string `json:"communication_with_lawyer" bson:"communication_with_lawyer"`

	// Есть ли у заключенных возможность направлять жалобы, ходатайства и заявления в надзирающие органы и правозащитные организации?
	CanPrisonersSubmitComplaints string `json:"can_prisoners_submit_complaints" bson:"can_prisoners_submit_complaints"`

	// Если ли у Вас есть дополнительная информация, которой Вы готовы поделиться с нами, то ее можно написать здесь:
	AdditionalInformation string `json:"additional_information" bson:"additional_information"`

	// Согласны ли Вы на публичную огласку приведенных Вами фактов
	PublicDisclosure string `json:"public_disclosure" bson:"public_disclosure"`

	// Готовы ли Вы сообщить свое имя и контакты?
	ProvideNameAndContacts string `json:"provide_name_and_contacts" bson:"provide_name_and_contacts"`

	// Согласны ли Вы на обработку Ваших персональных данных?
	ProcessingPersonalData string `json:"processing_personal_data" bson:"processing_personal_data"`

	// Какие нарушения, связанные с этапированием заключенных Вам известны?
	ViolationsStaging string `json:"violations_staging" bson:"violations_staging"`

	// С какими запретами (нарушениями) на отправление религиозных обрядов со стороны сотрудников ФСИН Вам приходилось сталкиваться?
	ViolationsReligiousRitesFromEmployees string `json:"violations_religious_rites_from_employees" bson:"violations_religious_rites_from_employees"`

	// С какими запретами (нарушениями) на отправление религиозных обрядов со стороны заключенных Вам приходилось сталкиваться?
	ViolationsReligiousRitesFromPrisoners string `json:"violations_religious_rites_from_prisoners" bson:"violations_religious_rites_from_prisoners"`

	// С какими нарушениями при применении мер взыскания, связанных с водворением в карцер и ШИЗО, переводом в ПКТ, ЕПКТ и на СУС, Вам приходилось сталкиваться?
	ViolationsPenaltiesRelatedToPlacement string `json:"violations_penalties_related_to_placement" bson:"violations_penalties_related_to_placement"`

	// Мы могли бы помочь Вам в составлении жалобы в Европейский суд по поводу нарушений в местах лишения свободы. Хотели бы Вы получить такую помощь?
	HelpEuropeanCourt string `json:"help_european_court" bson:"help_european_court"`

	// Источник поступления анкеты (если не Google Формы)
	Source string `json:"source"`

	// Одобрено для публикации?
	Approved string `json:"approved"`

	Position `json:"position"` // 32 и 33

	Warn string `json:"warn"` // 34 колонка

	ID      primitive.ObjectID `json:"_id" bson:"_id"`           // для монги
	PlaceID string             `json:"place_id" bson:"place_id"` // 35 колонка
}

type Position struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type PlaceCorona struct {
	Info        string `json:"info"`
	CommentFSIN string `json:"comment_fsin"`
	Date        string `json:"date"`

	Position `json:"position"`
	// широта - 4 колонка
	// долгота - 5 колонка
}

// структура учреждения ФСИН
type Place struct {
	ID primitive.ObjectID `bson:"_id" json:"_id"`

	// Полное название учреждния ФСИн
	Name string `json:"name"` // 0 колонка

	// Тип учреждения ФСИН
	Type string `json:"type"` // 1 колонка

	Location string `json:"location"` // 2 колонка

	// Доп информация from wiki
	Notes string `json:"notes"` // 3 колонка

	Position struct {
		Lat float64 `json:"lat"` // широта - 4 колонка
		Lng float64 `json:"lng"` // долгота - 5 колонка
	} `json:"position"`

	// Общее кол-во нарушений по нашей статистике
	NumberOfViolations uint64 `json:"number_of_violations" bson:"number_of_violations"` // 6 колонка

	// Номер телефона учреждения ФСИН
	Phones []string `json:"phones" bson:"phones"` // 7 колонка

	// время работы
	Hours string `json:"hours"` // 8 колонка

	Website string `json:"website"` // 9 колонка

	Address string `json:"address"` // 10 колонка

	Warn string `json:"warning"` // 11 колонка

	Coronavirus bool `json:"coronavirus" bson:"coronavirus"`

	Violations []Form `json:"violations" bson:"violations"`
}

// credentialsFile is the unmarshalled representation of a credentials file.
type CredentialsFile struct {
	Type string `json:"type"` // serviceAccountKey or userCredentialsKey

	// Service Account fields
	ClientEmail  string `json:"client_email"`
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	TokenURL     string `json:"token_uri"`
	ProjectID    string `json:"project_id"`

	// User Credential fields
	// (These typically come from gcloud auth.)
	ClientSecret string `json:"client_secret"`
	ClientID     string `json:"client_id"`
	RefreshToken string `json:"refresh_token"`
}

// прошлая структура для учреждения ФСИН
// Deprecated:
type OldPlace struct {
	// В каком регионе находится учреждение ФСИН о котором Вы рассказали?
	Region string `json:"region"`

	// О каком учреждении ФСИН Вы рассказали?
	FSINOrganization string `json:"fsin_organization"`

	// все сразу
	FullName string `json:"full_name"`
}
