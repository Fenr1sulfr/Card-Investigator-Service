package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type ConfirmationModel struct {
	DB *sql.DB
}

// Update to store result
func (c ConfirmationModel) Decline(cardID int64) error {
	query := `
		UPDATE cards
		SET status = 'decline'
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := c.DB.ExecContext(ctx, query, cardID)
	if err != nil {
		return err
	}
	return nil
}

func (c ConfirmationModel) Confirm(userID int64, cardID string) error {
	id, err := unFormatCardNumber(cardID)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO card_confirmations(card_id, user_id) VALUES ($1,$2)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = c.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}
	return nil
}

func (c ConfirmationModel) GetAllUsersByCard(cardID string) ([]*User, error) {
	id, err := unFormatCardNumber(cardID)
	if err != nil {
		return nil, err
	}
	query := `
		SELECT name,surname,email FROM card_confirmations
		LEFT JOIN users ON users.id=user_id
		WHERE card_id=$1;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}
	defer rows.Close()
	users := []*User{}
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Name,
			&user.Surname,
			&user.Email,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (c ConfirmationModel) GetAllCardsByUser(userEmail string) ([]*Card, error) {
	query := `
		SELECT id,case_number,
		registration_date,
		criminal_code_article,
		case_decision,case_summary,
		relation_to_event,
		invited_person_iin,
		invited_person_full_name,
		invited_person_position,
		organization_bin_or_iin,
		workplace,
		invited_person_status,
		planned_investigative_actions,
		scheduled_date_time,
		location,
		type_of_investigation,
		expected_outcome,
		investigator,
		is_business_related,
		pension_bin_or_iin,
		pension_workplace,
		enterpreneur_participation,
		defender_iin,
		defender_full_name
		FROM card_confirmations
LEFT JOIN cards ON cards.id=card_id
LEFT JOIN case_details ON cards.case_details_id = case_details.id
LEFT JOIN person_details ON cards.person_details_id = person_details.id
LEFT JOIN investigation_details ON cards.investigation_details_id = investigation_details.id
LEFT JOIN organizer_details ON cards.organizer_details_id = organizer_details.id
LEFT JOIN business_details ON cards.business_details_id = business_details.id
LEFT JOIN defender_details ON cards.defender_details_id = defender_details.id
INNER JOIN users ON card_confirmations.user_id = users.id
WHERE users.email = $1;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := c.DB.QueryContext(ctx, query, userEmail)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}
	defer rows.Close()
	cards := []*Card{}
	for rows.Next() {
		var card Card
		var RegIntNumberForString int
		err := rows.Scan(
			&RegIntNumberForString,
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
			return nil, err
		}
		card.BasicInfo.RegistrationNumber = formatCardNumber(RegIntNumberForString)
		cards = append(cards, &card)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return cards, nil
}
