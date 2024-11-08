package data

import "time"

type Card struct {
	RegistrationNumber          string    // Регистрационный номер (генерируется системой)
	CreationDate                time.Time // Дата создания документа (генерируется системой)
	CaseNumber                  string    // Номер УД (обязательное поле, ручной ввод, ФЛК 15 цифр)
	RegistrationDate            time.Time // Дата регистрации дела (автоподтягивание по номеру дела)
	CriminalCodeArticle         string    // Статья УК (автоподтягивание по номеру дела)
	CaseDecision                string    // Решение по делу (автоподтягивание по номеру дела)
	CaseSummary                 string    // Краткая фабула (автоподтягивание по номеру дела)
	InvitedPersonIIN            string    // ИИН вызываемого (обязательное поле, ручной ввод, ФЛК 12 цифр)
	InvitedPersonFullName       string    // ФИО вызываемого (автоподтягивание по ИИН вызываемого)
	InvitedPersonPosition       string    // Должность вызываемого (справочник)
	OrganizationBINOrIIN        string    // БИН/ИИН (обязательное поле, ручной ввод со стороны заполняющего, ФЛК 12 цифр)
	Workplace                   string    // Место работы (автоподтягивание по БИН/ИИН от заполняющего)
	Region                      string    // Регион (справочник)
	PlannedInvestigativeActions string    // Планируемые следственные действия (обязательное поле, ручной ввод)
	ScheduledDateTime           time.Time // Дата и время проведения (календарный и временной выбор)
	Location                    string    // Место проведения (справочник)
	Investigator                string    // Следователь (автоподтягивание с личного кабинета)
	InvitedPersonStatus         string    // Статус по делу вызываемого (справочник)
	RelationToEvent             string    // Отношение вызывающего к событию и субъекту (ручной ввод)
	TypeOfInvestigation         string    // Виды планируемого следствия (справочник)
	IsBusinessRelated           bool      // Относится ли к бизнесу (справочник)
	PensionBINOrIIN             string    // БИН/ИИН (пенсионные отчисления) - автоподтягивание последнего места работы
	PensionWorkplace            string    // Место работы (пенсионные отчисления) - автоподтягивание последнего места работы
	DefenderIIN                 string    // ИИН защитника (ручной ввод, ФЛК 12 цифр)
	DefenderFullName            string    // ФИО защитника (автоподтягивание по ИИН защитника)
	EntrepreneurParticipation   string    // Обоснование и необходимость участия предпринимателя (обязательное поле, ручной ввод)
	ExpectedOutcome             string    // Результат от планируемого следственного действия (обязательное поле, ручной ввод)}
}
