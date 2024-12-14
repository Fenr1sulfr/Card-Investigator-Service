package main

import (
	"api/internal/data"
	"api/internal/data/validator"
	"errors"
	"net/http"
	"time"
)

func (app *application) listCardsByRegion(w http.ResponseWriter, r *http.Request) {
	var input struct {
		data.Filters
		region string
	}
	v := validator.New()
	qs := r.URL.Query()
	input.region = app.readString(qs, "region", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	// if data.ValidateFilters(v, input.Filters); !v.Valid() {
	// 	app.failedValidationResponse(w, r, v.Errors)
	// 	return
	// }
	cards, err := app.models.Cards.GetAllByRegion(input.region)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorRespone(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"cards": cards}, nil)
	if err != nil {
		app.serverErrorRespone(w, r, err)
	}

}

func (app *application) getCard(w http.ResponseWriter, r *http.Request) {

	regnum := app.readRegistryNumberParam(r)

	//v := validator.New()
	//data.validateRegNum()
	//if !v.Valid();{app.failedValidationResponse(w,r,v.Errors) return}

	card, err := app.models.Cards.Get(regnum)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorRespone(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"card": card}, nil)
	if err != nil {
		app.serverErrorRespone(w, r, err)
	}
}

func (app *application) createCard(w http.ResponseWriter, r *http.Request) {
	var input struct {
		BasicInfo struct {
			RegistrationNumber string    `json:"registry_number"` // Регистрационный номер (генерируется системой)
			CreationDate       time.Time `json:"creation_date"`   // Дата создания документа (генерируется системой)
			Region             string    `json:"region"`          // Регион (справочник)
			Status             string    `json:"status"`
		} `json:"basic_info"`
		CaseDetails struct {
			CaseNumber          string    `json:"case_number"`       // Номер УД (обязательное поле, ручной ввод, ФЛК 15 цифр)
			RegistrationDate    time.Time `json:"registry_date"`     // Дата регистрации дела (автоподтягивание по номеру дела)
			CriminalCodeArticle string    `json:"criminal_core"`     // Статья УК (автоподтягивание по номеру дела)
			CaseDecision        string    `json:"case_decision"`     // Решение по делу (автоподтягивание по номеру дела)
			CaseSummary         string    `json:"case_summary"`      // Краткая фабула (автоподтягивание по номеру дела)
			RelationToEvent     string    `json:"relation_to_event"` // Отношение вызывающего к событию и субъекту (ручной ввод)
		} `json:"case_details"`
		PersonDetails struct {
			InvitedPersonIIN      string `json:"invited_person_iin"`       // ИИН вызываемого (обязательное поле, ручной ввод, ФЛК 12 цифр)
			InvitedPersonFullName string `json:"invited_person_full_name"` // ФИО вызываемого (автоподтягивание по ИИН вызываемого)
			InvitedPersonPosition string `json:"invited_person_position"`  // Должность вызываемого (справочник)
			OrganizationBINOrIIN  string `json:"organiztion_bin_or_iin"`   // БИН/ИИН (обязательное поле, ручной ввод со стороны заполняющего, ФЛК 12 цифр)
			Workplace             string `json:"workplace"`                // Место работы (автоподтягивание по БИН/ИИН от заполняющего)
			InvitedPersonStatus   string `json:"invited_person_status"`    // Статус по делу вызываемого (справочник)
		} `json:"person_details"`
		InvestigationDetails struct {
			PlannedInvestigativeActions string    `json:"planned_investigation_actions"` // Планируемые следственные действия (обязательное поле, ручной ввод)
			ScheduledDateTime           time.Time `json:"scheduled_date_time"`           // Дата и время проведения (календарный и временной выбор)
			Location                    string    `json:"location"`                      // Место проведения (справочник)
			TypeOfInvestigation         string    `json:"type_of_investigation"`         // Виды планируемого следствия (справочник)
			ExpectedOutcome             string    `json:"expected_outcome"`              // Результат от планируемого следственного действия (обязательное поле, ручной ввод)
		} `json:"investigation_details"`
		OrganizerDetails struct {
			Investigator string `json:"investigator"` // Следователь (автоподтягивание с личного кабинета)

		} `json:"organizer_details"`
		BusinessDetails struct {
			IsBusinessRelated         bool   `json:"is_business_related"`        // Относится ли к бизнесу (справочник)
			PensionBINOrIIN           string `json:"peron_bin_or_iin"`           // БИН/ИИН (пенсионные отчисления) - автоподтягивание последнего места работы
			PensionWorkplace          string `json:"pension_workplace"`          // Место работы (пенсионные отчисления) - автоподтягивание последнего места работы
			EntrepreneurParticipation string `json:"enterpreneur_participation"` // Обоснование и необходимость участия предпринимателя (обязательное поле, ручной ввод)
		} `json:"business_details"`
		DefenderDetails struct {
			DefenderIIN      string `json:"defender_iin"`       // ИИН защитника (ручной ввод, ФЛК 12 цифр)
			DefenderFullName string `json:"defender_full_name"` // ФИО защитника (автоподтягивание по ИИН защитника)
		} `json:"defender_details"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	card := &data.Card{
		BasicInfo:            input.BasicInfo,
		CaseDetails:          input.CaseDetails,
		PersonDetails:        input.PersonDetails,
		InvestigationDetails: input.InvestigationDetails,
		OrganizerDetails:     input.OrganizerDetails,
		BusinessDetails:      input.BusinessDetails,
		DefenderDetails:      input.DefenderDetails,
	}

	regNum, creationDate, err := app.models.Cards.Insert(*card)
	card.BasicInfo.RegistrationNumber = regNum
	card.BasicInfo.CreationDate = creationDate
	if err != nil {
		app.serverErrorRespone(w, r, err)
		return
	}
	// v := validator.New()
	// if data.ValidateCard(v, card); !v.Valid() {
	// 	app.failedValidationResponse(w, r, v.Errors)
	// 	return
	// }
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, data.ErrDuplicateEmail):
	// 		v.AddError("email", "a user with this email address already exists")
	// 		app.failedValidationResponse(w, r, v.Errors)
	// 	default:
	// 		app.serverErrorRespone(w, r, err)
	// 	}
	// 	return
	// }
	err = app.writeJSON(w, http.StatusCreated, envelope{"card": card}, nil)
	if err != nil {
		app.serverErrorRespone(w, r, err)
	}
}
