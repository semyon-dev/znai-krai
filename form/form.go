package form

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/semyon-dev/znai-krai/config"
	"net/http"
)

func Questions(c *gin.Context) {
	c.String(http.StatusOK, questions)
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

type question struct {
	Name     string `json:"name"`
	Question string `json:"question"`
	Required bool   `json:"required"`
}

var questions = ` { "questions":[
{"name": "status","question":"Какой Ваш статус?", "required":true},
{"name": "region","question":" В каком регионе находится учреждение ФСИН о котором Вы рассказали?", "required":true},
{"name": "fsin_organization","question":" О каком учреждении ФСИН Вы рассказали?", "required":true},
{"name": "time_of_offence","question":" Укажите когда произошли нарушения о которых Вы рассказали?"},
{"name": "physical_impact_from_employees","question":" С какими фактами применения физического воздействия со стороны сотрудников ФСИН Вам приходилось сталкиваться?"},
{"name": "physical_impact_from_prisoners","question":" С какими фактами применения физического воздействия со стороны заключенных Вам приходилось сталкиваться?"},
{"name": "psychological_impact_from_employees","question":" С какими фактами применения психологического воздействия со стороны отрудников ФСИН Вам приходилось сталкиваться?"},
{"name": "psychological_impact_from_prisoners","question":" С какими фактами применения психологического воздействия со стороны заключенных Вам приходилось сталкиваться?"},
{"name": "extortions_from_employees","question":" С какими фактами применения физического воздействия со стороны заключенных Вам приходилось сталкиваться?"},
{"name": "сorruption_from_employees","question":" Приходилось ли Вам сталкиваться с иными случаями коррупции сотрудников ФСИН?"},
{"name": "extortions_from_prisoners","question":" Приходилось ли Вам сталкиваться с фактами вымогательства со стороны заключенных?"},
{"name": "violations_of_medical_care","question":" Какие нарушения, связанные с оказанием медицинской помощи, Вы можете отметить?"},
{"name": "violations_of_food","question":" Какие нарушения, связанные с оказанием еды, Вы можете отметить?"},
{"name": "violations_of_clothes","question":" Какие нарушения, связанные с одеждой, Вы можете отметить?"},
{"name": "labor_slavery","question":" Известны ли Вам случаи трудового рабства?"},
{"name": "salary_of_prisoners","question":" Укажите, в каком диапазоне находится месячная зарплата заключенных?"},
{"name": "visits_with_relatives","question":" Какие нарушения, связанные с предоставлением свиданий с Родственниками, Вам известны?"},
{"name": "communication_with_lawyer","question":" Какие нарушения, связанные с общением с адвокатом (иным лицом, имеющим право на оказание юридической помощи), Вам известны?"},
{"name": "can_prisoners_submit_complaints","question":" Какие нарушения, связанные с общением с адвокатом (иным лицом, имеющим право на оказание юридической помощи), Вам известны?"},
{"name": "additional_information","question":" Если ли у Вас есть дополнительная информация, которой Вы готовы поделиться с нами, то ее можно написать здесь:"},
{"name": "public_disclosure","question":" Согласны ли Вы на публичную огласку приведенных Вами фактов?"},
{"name": "provide_name_and_contacts","question":" Готовы ли Вы сообщить свое имя и контакты?"},
{"name": "processing_personal_data","question":" Согласны ли Вы на обработку Ваших персональных данных?"},
{"name": "violations_staging","question":" Какие нарушения, связанные с этапированием заключенных Вам известны?"},
{"name": "violations_religious_rites_from_employees","question":" С какими запретами (нарушениями) на отправление религиозных обрядов со стороны сотрудников ФСИН Вам приходилось сталкиваться?"},
{"name": "violations_religious_rites_from_prisoners","question":" С какими запретами (нарушениями) на отправление религиозных обрядов со стороны заключенных Вам приходилось сталкиваться?"},
{"name": "violations_penalties_related_to_placement","question":" С какими нарушениями при применении мер взыскания, связанных с водворением в карцер и ШИЗО, переводом в ПКТ, ЕПКТ и на СУС, Вам приходилось сталкиваться?"},
{"name": "help_european_court","question":" Мы могли бы помочь Вам в составлении жалобы в Европейский суд по поводу нарушений в местах лишения свободы. Хотели бы Вы получить такую помощь?"}
]} `
