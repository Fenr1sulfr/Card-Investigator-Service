package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Card struct {
	BasicInfo            BasicInfo            `json:"basic_info"`
	CaseDetails          CaseDetails          `json:"case_details"`
	PersonDetails        PersonDetails        `json:"person_details"`
	InvestigationDetails InvestigationDetails `json:"investigation_details"`
	OrganizerDetails     OrganizerDetails     `json:"organizer_details"`
	BusinessDetails      BusinessDetails      `json:"business_details"`
	DefenderDetails      DefenderDetails      `json:"defender_details"`
}

// Basic information about the card.
type BasicInfo struct {
	RegistrationNumber string    `json:"registry_number"` // Регистрационный номер (генерируется системой)
	CreationDate       time.Time `json:"creation_date"`   // Дата создания документа (генерируется системой)
	Region             string    `json:"region"`          // Регион (справочник)
}

// Details about the specific case.
type CaseDetails struct {
	CaseNumber          string    `json:"case_number"`       // Номер УД (обязательное поле, ручной ввод, ФЛК 15 цифр)
	RegistrationDate    time.Time `json:"registry_date"`     // Дата регистрации дела (автоподтягивание по номеру дела)
	CriminalCodeArticle string    `json:"criminal_core"`     // Статья УК (автоподтягивание по номеру дела)
	CaseDecision        string    `json:"case_decision"`     // Решение по делу (автоподтягивание по номеру дела)
	CaseSummary         string    `json:"case_summary"`      // Краткая фабула (автоподтягивание по номеру дела)
	RelationToEvent     string    `json:"relation_to_event"` // Отношение вызывающего к событию и субъекту (ручной ввод)
}

// Information about the person being invited to the investigation.
type PersonDetails struct {
	InvitedPersonIIN      string `json:"invited_person_iin"`       // ИИН вызываемого (обязательное поле, ручной ввод, ФЛК 12 цифр)
	InvitedPersonFullName string `json:"invited_person_full_name"` // ФИО вызываемого (автоподтягивание по ИИН вызываемого)
	InvitedPersonPosition string `json:"invited_person_position"`  // Должность вызываемого (справочник)
	OrganizationBINOrIIN  string `json:"organiztion_bin_or_iin"`   // БИН/ИИН (обязательное поле, ручной ввод со стороны заполняющего, ФЛК 12 цифр)
	Workplace             string `json:"workplace"`                // Место работы (автоподтягивание по БИН/ИИН от заполняющего)
	InvitedPersonStatus   string `json:"invited_person_status"`    // Статус по делу вызываемого (справочник)
}

// Details specific to the planned investigation.
type InvestigationDetails struct {
	PlannedInvestigativeActions string    `json:"planned_investigation_actions"` // Планируемые следственные действия (обязательное поле, ручной ввод)
	ScheduledDateTime           time.Time `json:"scheduled_date_time"`           // Дата и время проведения (календарный и временной выбор)
	Location                    string    `json:"location"`                      // Место проведения (справочник)
	TypeOfInvestigation         string    `json:"type_of_investigation"`         // Виды планируемого следствия (справочник)
	ExpectedOutcome             string    `json:"expected_outcome"`              // Результат от планируемого следственного действия (обязательное поле, ручной ввод)
}

// Information about the organizer or investigator of the case.
type OrganizerDetails struct {
	Investigator string `json:"investigator"` // Следователь (автоподтягивание с личного кабинета)
}

// Business-related details of the investigation.
type BusinessDetails struct {
	IsBusinessRelated         bool   `json:"is_business_related"`        // Относится ли к бизнесу (справочник)
	PensionBINOrIIN           string `json:"peron_bin_or_iin"`           // БИН/ИИН (пенсионные отчисления) - автоподтягивание последнего места работы
	PensionWorkplace          string `json:"pension_workplace"`          // Место работы (пенсионные отчисления) - автоподтягивание последнего места работы
	EntrepreneurParticipation string `json:"enterpreneur_participation"` // Обоснование и необходимость участия предпринимателя (обязательное поле, ручной ввод)
}

// Information about the defender, if applicable.
type DefenderDetails struct {
	DefenderIIN      string `json:"defender_iin"`       // ИИН защитника (ручной ввод, ФЛК 12 цифр)
	DefenderFullName string `json:"defender_full_name"` // ФИО защитника (автоподтягивание по ИИН защитника)
}

type CardsModel struct {
	DB *sql.DB
}

var ErrNotRegNumber = errors.New("This is not a registry number")

// TODO:ADDING VALIDATOR
func unFormatCardNumber(cardNumber string) (int, error) {
	if len(cardNumber) < 3 && cardNumber[:2] != "Z-" {
		return -1, ErrNotRegNumber
	}
	atoi, err := strconv.Atoi(cardNumber[2:])
	if err != nil {
		return -1, err
	}
	return atoi, nil
}
func formatCardNumber(id int) string {
	// Определяем формат в зависимости от величины id
	if id < 1000 {
		return fmt.Sprintf("Z-%03d", id) // Добавляем ведущие нули для чисел меньше 1000
	}
	return fmt.Sprintf("Z-%d", id) // Для чисел 1000 и более без ведущих нулей
}

//TODO:ADDING VALIDATOR FOR sensitive data

func (m CardsModel) Insert(c Card) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := m.DB.BeginTx(ctx, nil)

	if err != nil {
		tx.Rollback()
		return err

	}
	var inputId struct {
		caseId          int64
		personId        int64
		investigationId int64
		organizerId     int64
		businessId      int64
		defenderId      int64
	}
	var query = `INSERT INTO case_details (case_number,registration_date,criminal_code_article,case_decision,case_summary,relation_to_event) VALUES
	($1,$2,$3,$4,$5,$6) RETURNING id`
	var args = []any{c.CaseDetails.CaseNumber, c.CaseDetails.RegistrationDate, c.CaseDetails.CriminalCodeArticle, c.CaseDetails.CaseDecision, c.CaseDetails.CaseDecision, c.CaseDetails.CaseSummary, c.CaseDetails.RelationToEvent}
	err = tx.QueryRowContext(ctx, query, args...).Scan(&inputId.caseId)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = `INSERT INTO perons_details (invited_person_iin,invited_person_full_name,invited_person_position,organization_bin_or_iin,workplace, invited_person_status) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`
	args = []any{c.PersonDetails.InvitedPersonIIN, c.PersonDetails.InvitedPersonFullName, c.PersonDetails.InvitedPersonPosition, c.BusinessDetails.PensionBINOrIIN, c.PersonDetails.Workplace, c.PersonDetails.InvitedPersonStatus}
	err = tx.QueryRowContext(ctx, query, args...).Scan(&inputId.personId)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = `INSERT INTO investigation_details (planned_investigative_actions,scheduled_date_time,location,type_of_investigation,expected_outcome) VALUES ($1,$2,$3,$4,$5) RETURNING id`
	args = []any{c.InvestigationDetails.PlannedInvestigativeActions, c.InvestigationDetails.ScheduledDateTime, c.InvestigationDetails.TypeOfInvestigation, c.InvestigationDetails.ExpectedOutcome}
	err = tx.QueryRowContext(ctx, query, args...).Scan(&inputId.investigationId)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = `INSERT INTO organizer_details (investigator) VALUES ($1) RETURNING id`
	args = []any{c.OrganizerDetails.Investigator}
	err = tx.QueryRowContext(ctx, query, args...).Scan(&inputId.organizerId)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = `INSERT INTO business_details (is_business_related, pension_bin_or_iin, pension_workplace,entrepreneur_participation) VALUES ($1,$2,$3,$4) RETURNING id`
	args = []any{c.BusinessDetails.IsBusinessRelated, c.BusinessDetails.PensionBINOrIIN, c.BusinessDetails.PensionWorkplace, c.BusinessDetails.EntrepreneurParticipation}
	err = tx.QueryRowContext(ctx, query, args...).Scan(&inputId.businessId)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = `INSERT INTO defender_details (defender_iin, defender_full_name) VALUES ($1,$2) RETURNING id`
	args = []any{c.DefenderDetails.DefenderIIN, c.DefenderDetails.DefenderFullName}
	err = tx.QueryRowContext(ctx, query, args...).Scan(&inputId.defenderId)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = `INSERT INTO cards (region,case_details_id,person_details_id,investigation_details_id,organizer_details_id,business_details_id,defender_details_id)`
	args = []any{c.BasicInfo.Region, inputId.caseId, inputId.personId, inputId.investigationId, inputId.organizerId, inputId.businessId, inputId.defenderId}
	err = tx.QueryRowContext(ctx, query, args...).Err()
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (m CardsModel) Get(regNumber string) (*Card, error) {
	id, err := unFormatCardNumber(regNumber)
	if err != nil {
		return nil, err
	}
	var card Card
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		SELECT  
    cards.creation_date, 
    cards.region, 
    case_details.case_number, 
    case_details.registration_date, 
    case_details.criminal_code_article, 
    case_details.case_decision, 
    case_details.case_summary, 
    case_details.relation_to_event, -
    person_details.invited_person_iin, 
    person_details.invited_person_full_name, 
    person_details.invited_person_position, 
    person_details.organization_bin_or_iin, 
    person_details.workplace, 
    person_details.invited_person_status, 
    investigation_details.planned_investigative_actions, 
    investigation_details.scheduled_date_time, 
    investigation_details.location, 
    investigation_details.type_of_investigation, 
    investigation_details.expected_outcome, 
    organizer_details.investigator, 
    business_details.is_business_related, 
    business_details.pension_bin_or_iin, 
    business_details.pension_workplace, 
    business_details.entrepreneur_participation, 
    defender_details.defender_iin, 
    defender_details.defender_full_name
FROM 
    cards
LEFT JOIN case_details ON cards.case_details_id = case_details.id
LEFT JOIN person_details ON cards.person_details_id = person_details.id
LEFT JOIN investigation_details ON cards.investigation_details_id = investigation_details.id
LEFT JOIN organizer_details ON cards.organizer_details_id = organizer_details.id
LEFT JOIN business_details ON cards.business_details_id = business_details.id
LEFT JOIN defender_details ON cards.defender_details_id = defender_details.id;
	WHERE cards.id=$1
	`
	err = m.DB.QueryRowContext(ctx, query, id).Scan(
		&card.BasicInfo.CreationDate,
		&card.BasicInfo.Region,
		&card.CaseDetails.CaseNumber,
		&card.CaseDetails.RegistrationDate,
		&card.CaseDetails.CriminalCodeArticle,
		&card.CaseDetails.CaseDecision,
		&card.CaseDetails.CaseSummary,
		&card.CaseDetails.RelationToEvent,
		&card.PersonDetails.InvitedPersonIIN,
		&card.PersonDetails.InvitedPersonFullName,
		&card.PersonDetails.InvitedPersonPosition,
		&card.PersonDetails.OrganizationBINOrIIN,
		&card.PersonDetails.Workplace,
		&card.PersonDetails.InvitedPersonStatus,
		&card.InvestigationDetails.PlannedInvestigativeActions,
		&card.InvestigationDetails.ScheduledDateTime,
		&card.InvestigationDetails.Location,
		&card.InvestigationDetails.TypeOfInvestigation,
		&card.InvestigationDetails.ExpectedOutcome,
		&card.OrganizerDetails.Investigator,
		&card.BusinessDetails.IsBusinessRelated,
		&card.BusinessDetails.PensionBINOrIIN,
		&card.BusinessDetails.PensionWorkplace,
		&card.BusinessDetails.EntrepreneurParticipation,
		&card.DefenderDetails.DefenderIIN,
		&card.DefenderDetails.DefenderFullName,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}
	return &card, nil
}