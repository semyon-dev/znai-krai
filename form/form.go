package form

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	log2 "github.com/semyon-dev/znai-krai/log"
	"github.com/semyon-dev/znai-krai/model"
	"net/http"
	"time"
)

func Questions(c *gin.Context) {

	var data = QuestionsData{}

	data.newQuestionTextfield(
		"region",
		"В каком регионе находится учреждение ФСИН о котором вы будете рассказывать?",
		true,
		"Новосибирск",
		"",
		nil,
	)

	data.newQuestionTextfield(
		"fsin_organization",
		"О каком учреждении ФСИН вы рассказываете?",
		true,
		"СИЗО 1",
		"",
		nil,
	)

	data.newQuestionTextfield(
		"time_of_offence",
		"Укажите когда, произошли нарушения, о которых вы хотите рассказать:",
		true,
		"2015-2018",
		"",
		nil,
	)

	data.newQuestionChooseOne(
		"status",
		"Какой ваш статус?",
		true,
		valuesStatusOther,
	)

	data.newQuestionChooseOne(
		"salary_of_prisoners",
		"Укажите, в каком диапазоне находится месячная зарплата заключенных?",
		false,
		model.ViolationsSalaryTypes,
	)

	data.newQuestionChooseOne(
		"violations_of_food_sure",
		"Были ли нарушения, связанные с питанием?",
		false,
		valuesYesNoDifficult,
	)

	data.newQuestionMultiply(
		"violations_of_food",
		"Какие нарушения, связанные с питанием, вы можете отметить?",
		false,
		"violations_of_food_sure",
		model.ViolationsFoodTypes,
	)

	data.newQuestionChooseOne(
		"violations_of_medical_care_sure",
		"Были ли нарушения, связанные с оказанием медицинской помощи?",
		false,
		valuesYesNoDifficult,
	)

	data.newQuestionMultiply(
		"violations_of_medical_care",
		"Какие нарушения, связанные с оказанием медицинской помощи, вы можете отметить?",
		false,
		"violations_of_medical_care_sure",
		model.ViolationsMedicalCareTypes,
	)

	data.newQuestionChooseOne(
		"physical_impact_sure",
		"Сталкивались ли вы с применением физического воздействия?",
		false,
		valuesYesNoDifficult,
	)

	data.newQuestionMultiply(
		"physical_impact_from_employees",
		"С какими фактами применения физического воздействия со стороны сотрудников ФСИН вам приходилось сталкиваться?",
		false,
		"physical_impact_sure",
		model.ViolationsPhysicalImpactTypes,
	)

	data.newQuestionMultiply(
		"physical_impact_from_prisoners",
		"С какими фактами применения физического воздействия со стороны заключенных вам приходилось сталкиваться?",
		false,
		"physical_impact_sure",
		model.ViolationsPhysicalImpactTypes,
	)

	data.newQuestionChooseOne(
		"communication_with_lawyer_sure",
		"Сталкивались ли вы с нарушениями, связанными с общением с адвокатом?",
		false,
		valuesYesNoDifficult,
	)

	data.newQuestionMultiply(
		"communication_with_lawyer",
		"Какие нарушения, связанные с общением с адвокатом (иным лицом, имеющим право на оказание юридической помощи), вам известны?",
		false,
		"communication_with_lawyer_sure",
		model.ViolationsCommunicationWithOthers,
	)

	data.newQuestionChooseOne(
		"visits_with_relatives_sure",
		"Сталкивались ли вы с нарушениями, связанными с предоставлением свиданий с родственниками?",
		false,
		valuesYesNoDifficult,
	)

	data.newQuestionMultiply(
		"visits_with_relatives",
		"Какие нарушения, связанные с предоставлением свиданий с родственниками, вам известны?",
		false,
		"visits_with_relatives_sure",
		model.ViolationsVisitsWithRelatives,
	)

	data.newQuestionChooseOne(
		"violations_of_clothes_sure",
		"Сталкивались ли вы с нарушениями, связанными с одеждой?",
		false,
		valuesYesNoDifficult,
	)

	data.newQuestionMultiply(
		"violations_of_clothes",
		"Какие нарушения, связанные с одеждой, вы можете отметить?",
		false,
		"violations_of_clothes_sure",
		model.ViolationsClothes,
	)

	data.newQuestionChooseOne(
		"violations_staging_sure",
		"Сталкивались ли вы с нарушениями, связанными с этапированием заключенных?",
		false,
		valuesYesNoDifficult,
	)

	data.newQuestionMultiply(
		"violations_staging",
		"Какие нарушения, связанные с этапированием заключенных вам известны?",
		false,
		"violations_staging_sure",
		model.ViolationsClothes,
	)

	data.newQuestionChooseOne(
		"violations_penalties_related_to_placement_sure",
		"Сталкивались ли вы с нарушениями при применении мер взыскания, связанных с водворением в карцер и ШИЗО, переводом в ПКТ, ЕПКТ и на СУС?",
		false,
		valuesYesNoDifficult,
	)

	data.newQuestionMultiply(
		"violations_penalties_related_to_placement",
		"С какими нарушениями при применении мер взыскания, связанных с водворением в карцер и ШИЗО, переводом в ПКТ, ЕПКТ и на СУС, вам приходилось сталкиваться?",
		false,
		"violations_penalties_related_to_placement_sure",
		model.ViolationsWithPlacementInPunishmentCellTypes,
	)

	data.newQuestionChooseOne(
		"violations_religious_rites_sure",
		"Сталкивались ли вы с нарушениями связанными с религией?",
		false,
		valuesYesNoDifficult,
	)

	data.newQuestionMultiply(
		"violations_religious_rites_from_employees",
		"С какими запретами (нарушениями) на отправление религиозных обрядов со стороны сотрудников ФСИН вам приходилось сталкиваться?",
		false,
		"violations_religious_rites_sure",
		model.ViolationsReligiousViolations,
	)

	data.newQuestionMultiply(
		"violations_religious_rites_from_prisoners",
		"С какими запретами (нарушениями) на отправление религиозных обрядов со стороны заключенных вам приходилось сталкиваться?",
		false,
		"violations_religious_rites_sure",
		model.ViolationsReligiousViolations,
	)

	data.newQuestionChooseOne(
		"psychological_impact_from_employees_sure",
		"Сталкивались ли вы с нарушениями, связанными с применением психологического воздействия?",
		false,
		valuesYesNoDifficult,
	)

	data.newQuestionMultiply(
		"psychological_impact_from_employees",
		"С какими фактами применения психологического воздействия со стороны сотрудников ФСИН вам приходилось сталкиваться?",
		false,
		"psychological_impact_from_employees_sure",
		model.ViolationsPsychologicalImpact,
	)

	data.newQuestionMultiply(
		"psychological_impact_from_prisoners",
		"С какими фактами применения психологического воздействия со стороны заключенных вам приходилось сталкиваться?",
		false,
		"psychological_impact_from_employees_sure",
		model.ViolationsPsychologicalImpact,
	)

	data.newQuestionChooseOne(
		"extortions_from_employees_sure",
		"Сталкивались ли вы с нарушениями,  связанными с фактами вымогательства со стороны сотрудников фсин?",
		false,
		valuesYesNoDifficult,
	)

	data.newQuestionMultiply(
		"extortions_from_employees",
		"В каких случаях вы сталкивались с фактами вымогательства со стороны сотрудников ФСИН?",
		false,
		"extortions_from_employees_sure",
		model.ViolationsExtortionsFromEmployeesTypes,
	)

	data.newQuestionChooseOne(
		"corruption_from_employees",
		"Приходилось ли вам сталкиваться с иными случаями коррупции сотрудников ФСИН?",
		false,
		valuesYesNoDifficultOther,
	)

	data.newQuestionChooseOne(
		"extortions_from_prisoners",
		"Приходилось ли вам сталкиваться с фактами вымогательства со стороны заключенных?",
		false,
		valuesYesNoDifficultOther,
	)

	data.newQuestionChooseOne(
		"can_prisoners_submit_complaints",
		"Есть ли у заключенных возможность направлять жалобы, ходатайства и заявления в надзирающие органы и правозащитные организации?",
		false,
		valuesYesNoDifficultOther,
	)

	data.newQuestionChooseOne(
		"labor_slavery",
		"Известны ли вам случаи трудового рабства?",
		false,
		valuesYesNoDifficultOther,
	)

	data.newQuestionChooseOne(
		"help_european_court",
		"Мы могли бы помочь вам в составлении жалобы в Европейский суд по поводу нарушений в местах лишения свободы. Хотели бы вы получить такую помощь?",
		false,
		valuesYesNoOther,
	)

	data.newQuestionChooseOne(
		"public_disclosure",
		"Согласны ли Вы на публичную огласку приведенных вами фактов?",
		true,
		valuesYesNo,
	)

	data.newQuestionChooseOne(
		"processing_personal_data",
		"Согласны ли вы на обработку Ваших персональных данных?",
		true,
		valuesYes,
	)

	data.newQuestionTextarea(
		"additional_information",
		"Если ли у вас есть дополнительная информация, которой вы готовы поделиться с нами, то ее можно написать здесь:",
		false,
		nil,
	)

	data.newQuestionChooseOne(
		"provide_name_and_contacts_sure",
		"Готовы ли вы сообщить свое имя и контакты?",
		false,
		valuesYesNo,
	)

	data.newQuestionTextfield(
		"contacts",
		"Ваши контакты:",
		false,
		"Мамонтов Власий Демьянович, 89001112233",
		"provide_name_and_contacts_sure",
		nil,
	)

	data.newQuestionWithHtml(
		"add_files",
		"Если ли у вас есть файлы которые относятся к нарушениям, то можете загрузить их здесь:",
		"<iframe width=\"250\" height=\"54\" frameborder=\"0\" src=\"https://mega.nz/drop#!0SWpxKkiXk4!l!en\"></iframe>",
		false,
	)

	c.JSON(http.StatusOK, data)
}

func (data *QuestionsData) newQuestion(name, questionName, typ string, required bool, hint, requires string, values []string, html string) {
	*data = append(*data, question{
		Name:     name,
		Question: questionName,
		Required: required,
		Type:     typ,
		Values:   values,
		Hint:     hint,
		Requires: requires,
		Html:     html,
	})
}

func (data *QuestionsData) newQuestionMultiply(name, questionName string, required bool, requires string, values []string) {
	data.newQuestion(name, questionName, "choose_multiply", required, "", requires, values, "")
}

func (data *QuestionsData) newQuestionChooseOne(name, questionName string, required bool, values []string) {
	data.newQuestion(name, questionName, "choose_one", required, "", "", values, "")
}

func (data *QuestionsData) newQuestionTextfield(name, questionName string, required bool, hint, requires string, values []string) {
	data.newQuestion(name, questionName, "textfield", required, hint, requires, values, "")
}

func (data *QuestionsData) newQuestionTextarea(name, questionName string, required bool, values []string) {
	data.newQuestion(name, questionName, "textarea", required, "", "", values, "")
}

// for upload files
func (data *QuestionsData) newQuestionWithHtml(name, questionName, html string, required bool) {
	data.newQuestion(name, questionName, "html", required, "", "", nil, html)
}

var valuesYes = []string{"Да"}
var valuesYesNo = []string{"Да", "Нет"}
var valuesYesNoOther = []string{"Да", "Нет", "другое"}
var valuesYesNoDifficult = []string{"Да", "Нет", "Затрудняюсь ответить"}
var valuesYesNoDifficultOther = []string{"Да", "Нет", "Затрудняюсь ответить", "другое"}
var valuesStatusOther = []string{"Бывший заключенный", "Родственник заключенного", "Заключенный в настоящее время", "Адвокат", "другое"}

// для кнопки "сообщить об ошибке"
type report struct {
	Email      string `json:"email"` // почта для обратной связи
	Bug        string `json:"bug"`
	PlaceId    string `json:"place_id"`
	NameOfFSIN string `json:"fsin_organization"`
}

func Report(c *gin.Context) {
	var report report
	err := c.ShouldBind(&report)
	if err != nil {
		log.Error().Err(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "вad request" + err.Error(),
		})
		return
	}
	hooked := log.Hook(log2.ReportHook{})
	hooked.Error().Msg("Новый report:\n" + "Текст: " + report.Bug + "\nemail: " + report.Email + "\nНазвание МЛС: " + report.NameOfFSIN + "\nplace_id: " + report.PlaceId + "\n" + "Время: " + time.Now().Format("2006.01.02 15:04:05") + "\nOrigin: " + c.GetHeader("Origin") + "\nHost: " + c.Request.Host + "\nClientIP: " + c.ClientIP())
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

type QuestionsData []question

type question struct {
	Name     string   `json:"name"`
	Question string   `json:"question"`
	Required bool     `json:"required"`
	Requires string   `json:"requires"`
	Type     string   `json:"type"`
	Values   []string `json:"values"`
	Hint     string   `json:"hint"`
	Html     string   `json:"button"`
}
