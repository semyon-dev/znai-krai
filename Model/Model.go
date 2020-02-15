package Model

// структура для нарушений
type Form struct {
	// отметка времени
	Time string `json:"time" binding:"required"`
	// Какой Ваш статус?
	Status string `json:"status" binding:"required"`
	// В каком регионе находится учреждение ФСИН о котором Вы рассказали?
	Region string `json:"region" binding:"required"`
	// О каком учреждении ФСИН Вы рассказали?
	FSINОrganization string `json:"fsin_organization" binding:"required"`
	// Укажите когда произошли нарушения о которых Вы рассказали?
	TimeOfOffence string `json:"time_of_offence" binding:"required"`
}
