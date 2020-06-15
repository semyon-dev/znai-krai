package form

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/model"
	"net/http"
)

func Questions(c *gin.Context) {

	var data QuestionsData

	var region = question{
		Name:     "region",
		Question: "В каком регионе находится учреждение ФСИН о котором вы будете рассказывать?",
		Required: true,
		Type:     "textfield",
		Hint:     "Новосибирск",
	}

	var fsinOrganization = question{
		Name:     "fsin_organization",
		Question: "О каком учреждении ФСИН вы рассказываете?",
		Required: true,
		Type:     "textfield",
		Hint:     "СИЗО 1",
	}

	var timeOfOffence = question{
		Name:     "time_of_offence",
		Question: "Укажите когда произошли нарушения о которых вы хотите рассказать?",
		Required: true,
		Type:     "textfield",
		Hint:     "2015-2018",
	}

	var questionStatus = question{
		Name:     "status",
		Question: "Какой ваш статус?",
		Required: true,
		Type:     "choose_one",
		Values:   []string{"Бывший заключенный", "Родственник заключенного", "Заключенный в настоящее время", "Адвокат", "другое"},
	}

	var salaryOfPrisoners = question{
		Name:     "salary_of_prisoners",
		Question: "Укажите, в каком диапазоне находится месячная зарплата заключенных?",
		Required: false,
		Type:     "choose_one",
		Values:   model.ViolationsSalaryTypes,
	}

	var violationFood = question{
		Name:     "violations_of_food",
		Question: "Какие нарушения, связанные с питанием, вы можете отметить?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsFoodTypes,
	}

	var violationsMedicalCare = question{
		Name:     "violations_of_medical_care",
		Question: "Какие нарушения, связанные с оказанием медицинской помощи, вы можете отметить?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsMedicalCareTypes,
	}

	var contacts = question{
		Name:     "provide_name_and_contacts",
		Question: "Готовы ли вы сообщить свое имя и контакты? Если нет - пропустите вопрос.",
		Required: false,
		Type:     "textfield",
		Hint:     "Мамонтов Власий Демьянович, 89001112233",
	}

	var physicalImpactFromEmployees = question{
		Name:     "physical_impact_from_employees",
		Question: "С какими фактами применения физического воздействия со стороны сотрудников ФСИН вам приходилось сталкиваться?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsPhysicalImpactTypes,
	}

	var physicalImpactFromPrisoners = question{
		Name:     "physical_impact_from_prisoners",
		Question: "С какими фактами применения физического воздействия со стороны заключенных вам приходилось сталкиваться?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsPhysicalImpactTypes,
	}

	var visitsWithRelatives = question{
		Name:     "visits_with_relatives",
		Question: "Какие нарушения, связанные с предоставлением свиданий с Родственниками, вам известны?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsVisitsWithRelatives,
	}

	var violationsOfClothes = question{
		Name:     "violations_of_clothes",
		Question: "Какие нарушения, связанные с одеждой, вы можете отметить?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsClothes,
	}

	var violationsStaging = question{
		Name:     "violations_staging",
		Question: "Какие нарушения, связанные с этапированием заключенных вам известны?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsClothes,
	}

	var violationsPenaltiesRelatedToPlacement = question{
		Name:     "violations_penalties_related_to_placement",
		Question: "С какими нарушениями при применении мер взыскания, связанных с водворением в карцер и ШИЗО, переводом в ПКТ, ЕПКТ и на СУС, вам приходилось сталкиваться?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsWithPlacementInPunishmentCellTypes,
	}

	var communicationWithLawyer = question{
		Name:     "communication_with_lawyer",
		Question: "Какие нарушения, связанные с общением с адвокатом (иным лицом, имеющим право на оказание юридической помощи), вам известны?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsCommunicationWithOthers,
	}

	var violationsReligiousRitesFromEmployees = question{
		Name:     "violations_religious_rites_from_employees",
		Question: "С какими запретами (нарушениями) на отправление религиозных обрядов со стороны сотрудников ФСИН вам приходилось сталкиваться?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsReligiousViolations,
	}
	var violationsReligiousRitesFromPrisoners = question{
		Name:     "violations_religious_rites_from_prisoners",
		Question: "С какими запретами (нарушениями) на отправление религиозных обрядов со стороны заключенных вам приходилось сталкиваться?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsReligiousViolations,
	}

	var psychologicalImpactFromEmployees = question{
		Name:     "psychological_impact_from_employees",
		Question: "С какими фактами применения психологического воздействия со стороны сотрудников ФСИН вам приходилось сталкиваться?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsPsychologicalImpact,
	}

	var psychologicalImpactFromPrisoners = question{
		Name:     "psychological_impact_from_prisoners",
		Question: "С какими фактами применения психологического воздействия со стороны заключенных вам приходилось сталкиваться?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsPsychologicalImpact,
	}

	var extortionsFromEmployees = question{
		Name:     "extortions_from_employees",
		Question: "В каких случаях вы сталкивались с фактами вымогательства со стороны сотрудников ФСИН?",
		Required: false,
		Type:     "choose_multiply",
		Values:   model.ViolationsExtortionsFromEmployeesTypes,
	}

	var corruptionFromEmployees = question{
		Name:     "corruption_from_employees",
		Question: "Приходилось ли вам сталкиваться с иными случаями коррупции сотрудников ФСИН?",
		Required: false,
		Type:     "choose_one",
		Values:   []string{"Да", "Нет", "Затрудняюсь ответить", "другое"},
	}

	var extortionsFromPrisoners = question{
		Name:     "extortions_from_prisoners",
		Question: "Приходилось ли вам сталкиваться с фактами вымогательства со стороны заключенных?",
		Required: false,
		Type:     "choose_one",
		Values:   []string{"Да", "Нет", "Затрудняюсь ответить", "другое"},
	}

	var canPrisonersSubmitComplaints = question{
		Name:     "can_prisoners_submit_complaints",
		Question: "Есть ли у заключенных возможность направлять жалобы, ходатайства и заявления в надзирающие органы и правозащитные организации?",
		Required: false,
		Type:     "choose_one",
		Values:   []string{"Да", "Нет", "Затрудняюсь ответить", "другое"},
	}

	var laborSlavery = question{
		Name:     "labor_slavery",
		Question: "Известны ли вам случаи трудового рабства?",
		Required: false,
		Type:     "choose_one",
		Values:   []string{"Да", "Нет", "Затрудняюсь ответить", "другое"},
	}

	var helpEuropeanCourt = question{
		Name:     "help_european_court",
		Question: "Мы могли бы помочь вам в составлении жалобы в Европейский суд по поводу нарушений в местах лишения свободы. Хотели бы вы получить такую помощь?",
		Required: false,
		Type:     "choose_one",
		Values:   []string{"Да", "Нет", "другое"},
	}

	var questionPublicDisclosure = question{
		Name:     "public_disclosure",
		Question: "Согласны ли Вы на публичную огласку приведенных вами фактов?",
		Required: true,
		Type:     "choose_one",
		Values:   []string{"Да", "Нет"},
	}

	var questionProcessingPersonalData = question{
		Name:     "processing_personal_data",
		Question: "Согласны ли вы на обработку Ваших персональных данных?",
		Required: true,
		Type:     "choose_one",
		Values:   []string{"Да"},
	}

	var additionalInformation = question{
		Name:     "additional_information",
		Question: "Если ли у вас есть дополнительная информация, которой вы готовы поделиться с нами, то ее можно написать здесь:",
		Required: false,
		Type:     "textarea",
	}

	var addFiles = question{
		Name:     "add_files",
		Question: "Если ли у вас есть файлы которые относятся к нарушениям, то можете загрузить их здесь:",
		Required: false,
		Type:     "button",
		Button:   "<iframe width=\"250\" height=\"54\" frameborder=\"0\" src=\"https://mega.nz/drop#!0SWpxKkiXk4!d!en\"></iframe>",
	}

	data = append(
		data,
		questionStatus,
		region,
		timeOfOffence,
		fsinOrganization,
		contacts,
		physicalImpactFromEmployees,
		physicalImpactFromPrisoners,
		extortionsFromEmployees,
		corruptionFromEmployees,
		extortionsFromPrisoners,
		psychologicalImpactFromPrisoners,
		psychologicalImpactFromEmployees,
		violationsReligiousRitesFromEmployees,
		violationsReligiousRitesFromPrisoners,
		communicationWithLawyer,
		visitsWithRelatives,
		canPrisonersSubmitComplaints,
		violationsOfClothes,
		laborSlavery,
		salaryOfPrisoners,
		violationFood,
		violationsMedicalCare,
		violationsStaging,
		violationsPenaltiesRelatedToPlacement,
		additionalInformation,
		helpEuropeanCourt,
		questionPublicDisclosure,
		questionProcessingPersonalData,
		addFiles,
	)

	c.JSON(http.StatusOK, data)
}

func Report(c *gin.Context) {
	var report report
	err := c.ShouldBind(&report)
	if err != nil {
		log.Error().Err(err)
		c.JSON(http.StatusBadRequest, report)
		return
	}
	hooked := log.Hook(config.ReportHook{})
	hooked.Error().Msg("new report: " + report.Bug + ", email: " + report.Email)
	c.JSON(http.StatusOK, report)
}

// для кнопки "сообщить об ошибке"
type report struct {
	Email string `json:"email"` // почта для обратной связи
	Bug   string `json:"bug"`
}

type QuestionsData []question

type question struct {
	Name     string   `json:"name"`
	Question string   `json:"question"`
	Required bool     `json:"required"`
	Type     string   `json:"type"`
	Values   []string `json:"values"`
	Hint     string   `json:"hint"`
	Button   string   `json:"button"`
}
