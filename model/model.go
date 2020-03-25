package model

// структура для нарушений
type Form struct {
	// отметка времени
	Time string `json:"time"`
	// Какой Ваш статус?
	Status string `json:"status" binding:"required"`
	// В каком регионе находится учреждение ФСИН о котором Вы рассказали?
	Region string `json:"region" binding:"required"`
	// О каком учреждении ФСИН Вы рассказали?
	FSINОrganization string `json:"fsin_organization" binding:"required"`
	// Укажите когда произошли нарушения о которых Вы рассказали?
	TimeOfOffence string `json:"time_of_offence" `
	// С какими фактами применения физического воздействия со стороны сотрудников ФСИН Вам приходилось сталкиваться?
	PhysicalImpactFromEmployees string `json:"physical_impact_from_employees"`
	// С какими фактами применения физического воздействия со стороны заключенных Вам приходилось сталкиваться?
	PhysicalImpactFromPrisoners string `json:"physical_impact_from_prisoners"`

	// С какими фактами психологического воздействия со стороны сотрудников ФСИН Вам приходилось сталкиваться?
	PsychologicalImpactFromEmployees string `json:"psychological_impact_from_employees"`
	// С какими фактами психологического воздействия со стороны заключенных Вам приходилось сталкиваться?
	PsychologicalImpactFromPrisoners string `json:"psychological_impact_from_prisoners"`

	// В каких случаях Вы сталкивались с фактами вымогательства со стороны сотрудников ФСИН?
	ExtortionsFromEmployees string `json:"extortions_from_employees"`

	// Приходилось ли Вам сталкиваться с иными случаями коррупции сотрудников ФСИН?
	СorruptionFromEmployees string `json:"сorruption_from_employees"`

	// Приходилось ли Вам сталкиваться с фактами вымогательства со стороны заключенных?
	ExtortionsFromPrisoners string `json:"extortions_from_prisoners"`

	// Какие нарушения, связанные с оказанием медицинской помощи, Вы можете отметить?
	ViolationsOfMedicalCare string `json:"violations_of_medical_care"`

	// Какие нарушения, связанные с питанием, Вы можете отметить?
	ViolationsOfFood string `json:"violations_of_food"`

	// Какие нарушения, связанные с одеждой, Вы можете отметить?
	ViolationsOfClothes string `json:"violations_of_clothes"`

	// Известны ли Вам случаи трудового рабства?
	LaborSlavery string `json:"labor_slavery"`

	// Укажите, в каком диапазоне находится месячная зарплата заключенных?
	SalaryOfPrisoners string `json:"salary_of_prisoners"`

	// Какие нарушения, связанные с предоставлением свиданий с Родственниками, Вам известны?
	VisitsWithRelatives string `json:"visits_with_relatives"`

	// Какие нарушения, связанные с иными формами общения с Родственниками, Вам известны?
	CommunicationWithRelatives string `json:"communication_with_relatives"`

	// Какие нарушения, связанные с общением с адвокатом (иным лицом, имеющим право на оказание юридической помощи), Вам известны?
	CommunicationWithLawyer string `json:"communication_with_lawyer"`

	// Есть ли у заключенных возможность направлять жалобы, ходатайства и заявления в надзирающие органы и правозащитные организации?
	CanPrisonersSubmitComplaints string `json:"can_prisoners_submit_complaints"`

	// Если ли у Вас есть дополнительная информация, которой Вы готовы поделиться с нами, то ее можно написать здесь:
	AdditionalInformation string `json:"additional_information"`

	// Согласны ли Вы на публичную огласку приведенных Вами фактов?
	PublicDisclosure string `json:"public_disclosure"`

	// Готовы ли Вы сообщить свое имя и контакты?
	ProvideNameAndContacts string `json:"provide_name_and_contacts"`

	// Согласны ли Вы на обработку Ваших персональных данных?
	ProcessingPersonalData string `json:"processing_personal_data"`

	// Какие нарушения, связанные с этапированием заключенных Вам известны?
	ViolationsStaging string `json:"violations_staging"`

	// С какими запретами (нарушениями) на отправление религиозных обрядов со стороны сотрудников ФСИН Вам приходилось сталкиваться?
	ViolationsReligiousRitesFromEmployees string `json:"violations_religious_rites_from_employees"`

	// С какими запретами (нарушениями) на отправление религиозных обрядов со стороны заключенных Вам приходилось сталкиваться?
	ViolationsReligiousRitesFromPrisoners string `json:"violations_religious_rites_from_prisoners"`

	// С какими нарушениями при применении мер взыскания, связанных с водворением в карцер и ШИЗО, переводом в ПКТ, ЕПКТ и на СУС, Вам приходилось сталкиваться?
	ViolationsPenaltiesRelatedToPlacement string `json:"violations_penalties_related_to_placement"`

	// Мы могли бы помочь Вам в составлении жалобы в Европейский суд по поводу нарушений в местах лишения свободы. Хотели бы Вы получить такую помощь?
	HelpEuropeanCourt string `json:"help_european_court"`

	// Источник поступления анкеты (если не Google Формы)
	Source string `json:"source"`

	// Одобрено для публикации?
	//Approved bool
}

// структура учреждения ФСИН
type Place struct {

	// Полное название учреждния ФСИн
	Name string `json:"name"`

	// Тип учреждения ФСИН
	Type string `json:"type"`

	Location string `json:"location"`

	// Доп информация
	Notes string `json:"notes"`

	Position struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"position"`

	// Общее кол-во нарушений по нашей статистике
	NumberOfViolations uint64 `json:"number_of_violations"`

	// Номер телефона учреждения ФСИН
	PhoneNumber string `json:"phone_number"`

	// Google Maps Rating
	GoogleMapsRating float64 `json:"google_maps_rating"`

	Website string `json:"website"`
}

// Deprecated: структура для учреждения ФСИН
type OldPlace struct {
	// В каком регионе находится учреждение ФСИН о котором Вы рассказали?
	Region string `json:"region"`

	// О каком учреждении ФСИН Вы рассказали?
	FSINОrganization string `json:"fsin_organization"`

	// все сразу
	FullName string `json:"full_name"`
}
