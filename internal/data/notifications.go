package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Notification struct {
	CaseNumber      string    `json:"case_number"`
	RecipientName   string    `json:"invited_person_full_name"`
	Date            time.Time `json:"scheduled_date_time"`
	Location        string    `json:"location"`
	CodeArticle     string    `json:"criminal_code_article"`
	Investigator    string    `json:"investigator"`
	RelationToEvent string    `json:"relation_to_event"`
}

type NotificationModel struct {
	DB *sql.DB
}

func (n NotificationModel) GetNotificationInfo(regNumber string) (*Notification, error) {
	id, err := unFormatCardNumber(regNumber)
	if err != nil {
		return nil, err
	}
	query := `
	SELECT 
    case_details.case_number,
    person_details.invited_person_full_name,
    investigation_details.scheduled_date_time,
    investigation_details.location,
    case_details.criminal_code_article,
    organizer_details.investigator,
    case_details.relation_to_event
FROM 
    cards
JOIN 
    case_details ON cards.case_details_id = case_details.id
JOIN 
    person_details ON cards.person_details_id = person_details.id
JOIN 
    investigation_details ON cards.investigation_details_id = investigation_details.id
JOIN 
    organizer_details ON cards.organizer_details_id = organizer_details.id
	WHERE cards.id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var notification Notification

	err = n.DB.QueryRowContext(ctx, query, id).Scan(
		&notification.CaseNumber,
		&notification.RecipientName,
		&notification.Date,
		&notification.Location,
		&notification.CodeArticle,
		&notification.Investigator,
		&notification.RelationToEvent,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}
	return &notification, nil
}
