package form

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/semyon-dev/znai-krai/config"
	"github.com/semyon-dev/znai-krai/model"
	"net/http"
)

func Questions(c *gin.Context) {
	var data QuestionsData
	err := json.Unmarshal([]byte(questions), &data)
	if err != nil {
		fmt.Println(err.Error())
	}
	for i, v := range data {
		if v.Type == "" {
			data[i].Type = "textarea"
		}
	}

	var questionStatus = question{
		Name:     "status",
		Question: "Какой ваш статус?",
		Required: true,
		Type:     "choose_one",
		Values:   []string{"Бывший заключенный", "Родственник заключенного", "Заключенный в настоящее время", "Адвокат", "другое"},
	}

	var questionPublicDisclosure = question{
		Name:     "public_disclosure",
		Question: "Согласны ли Вы на публичную огласку приведенных Вами фактов?",
		Required: true,
		Type:     "choose_one",
		Values:   []string{"Да", "Нет"},
	}

	var questionProcessingPersonalData = question{
		Name:     "processing_personal_data",
		Question: "Согласны ли Вы на обработку Ваших персональных данных?",
		Required: true,
		Type:     "choose_one",
		Values:   []string{"Да"},
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
		Question: "Какие нарушения, связанные с оказанием еды, вы можете отметить?",
		Required: false,
		Type:     "choose_one",
		Values:   model.ViolationsFoodTypes,
	}

	var violationsMedicalCare = question{
		Name:     "violations_of_medical_care",
		Question: "Какие нарушения, связанные с оказанием медицинской помощи, вы можете отметить?",
		Required: false,
		Type:     "choose_one",
		Values:   model.ViolationsMedicalCareTypes,
	}

	var contacts = question{
		Name:     "provide_name_and_contacts",
		Question: "Готовы ли вы сообщить свое имя и контакты? Если нет - пропустите поле.",
		Required: false,
		Type:     "textfield",
	}

	var physicalImpactFromEmployees = question{
		Name:     "physical_impact_from_employees",
		Question: "С какими фактами применения физического воздействия со стороны сотрудников ФСИН Вам приходилось сталкиваться?",
		Required: false,
		Type:     "choose_one",
		Values:   model.ViolationsPhysicalImpactTypes,
	}

	var physicalImpactFromPrisoners = question{
		Name:     "physical_impact_from_prisoners",
		Question: "С какими фактами применения физического воздействия со стороны заключенных вам приходилось сталкиваться?",
		Required: false,
		Type:     "choose_one",
		Values:   model.ViolationsPhysicalImpactTypes,
	}

	data = append(data, questionStatus, questionPublicDisclosure, questionProcessingPersonalData, salaryOfPrisoners, violationFood, violationsMedicalCare, contacts, physicalImpactFromEmployees, physicalImpactFromPrisoners)
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

// TODO: убрать потом
var types = [...]string{
	"text",
	"choose_one",
	"choose_multiply",
}

type question struct {
	Name     string   `json:"name"`
	Question string   `json:"question"`
	Required bool     `json:"required"`
	Type     string   `json:"type"`
	Values   []string `json:"values"`
}

var questions = ` [
{"name": "region","question":" В каком регионе находится учреждение ФСИН о котором Вы рассказали?", "required":true},
{"name": "fsin_organization","question":" О каком учреждении ФСИН Вы рассказали?", "required":true},
{"name": "time_of_offence","question":" Укажите когда произошли нарушения о которых Вы рассказали?"},
{"name": "psychological_impact_from_employees","question":" С какими фактами применения психологического воздействия со стороны отрудников ФСИН Вам приходилось сталкиваться?"},
{"name": "psychological_impact_from_prisoners","question":" С какими фактами применения психологического воздействия со стороны заключенных Вам приходилось сталкиваться?"},
{"name": "extortions_from_employees","question":" С какими фактами применения физического воздействия со стороны заключенных Вам приходилось сталкиваться?"},
{"name": "сorruption_from_employees","question":" Приходилось ли Вам сталкиваться с иными случаями коррупции сотрудников ФСИН?"},
{"name": "extortions_from_prisoners","question":" Приходилось ли Вам сталкиваться с фактами вымогательства со стороны заключенных?"},
{"name": "violations_of_clothes","question":" Какие нарушения, связанные с одеждой, Вы можете отметить?"},
{"name": "labor_slavery","question":" Известны ли Вам случаи трудового рабства?"},
{"name": "visits_with_relatives","question":" Какие нарушения, связанные с предоставлением свиданий с Родственниками, Вам известны?"},
{"name": "communication_with_lawyer","question":" Какие нарушения, связанные с общением с адвокатом (иным лицом, имеющим право на оказание юридической помощи), Вам известны?"},
{"name": "can_prisoners_submit_complaints","question":" Какие нарушения, связанные с общением с адвокатом (иным лицом, имеющим право на оказание юридической помощи), Вам известны?"},
{"name": "additional_information","question":" Если ли у Вас есть дополнительная информация, которой Вы готовы поделиться с нами, то ее можно написать здесь:"},
{"name": "violations_staging","question":" Какие нарушения, связанные с этапированием заключенных Вам известны?"},
{"name": "violations_religious_rites_from_employees","question":" С какими запретами (нарушениями) на отправление религиозных обрядов со стороны сотрудников ФСИН Вам приходилось сталкиваться?"},
{"name": "violations_religious_rites_from_prisoners","question":" С какими запретами (нарушениями) на отправление религиозных обрядов со стороны заключенных Вам приходилось сталкиваться?"},
{"name": "violations_penalties_related_to_placement","question":" С какими нарушениями при применении мер взыскания, связанных с водворением в карцер и ШИЗО, переводом в ПКТ, ЕПКТ и на СУС, Вам приходилось сталкиваться?"},
{"name": "help_european_court","question":" Мы могли бы помочь Вам в составлении жалобы в Европейский суд по поводу нарушений в местах лишения свободы. Хотели бы Вы получить такую помощь?"}
]`
