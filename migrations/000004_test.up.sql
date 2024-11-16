
CREATE TABLE case_details (
    id SERIAL PRIMARY KEY,
    case_number VARCHAR(15) NOT NULL CHECK (LENGTH(case_number) = 15), -- Case number, 15-digit format
    registration_date TIMESTAMP NOT NULL,          -- Registration date of the case
    criminal_code_article VARCHAR(50),             -- Criminal code article
    case_decision TEXT,                            -- Decision on the case
    case_summary TEXT,                             -- Case summary
    relation_to_event TEXT                         -- Relationship of the caller to the event
);

CREATE TABLE person_details (
    id SERIAL PRIMARY KEY,
    invited_person_iin VARCHAR(12) NOT NULL CHECK (LENGTH(invited_person_iin) = 12), -- IIN, 12 digits
    invited_person_full_name VARCHAR(255),    -- Full name of the invited person
    invited_person_position VARCHAR(100),     -- Position
    organization_bin_or_iin VARCHAR(12) NOT NULL CHECK (LENGTH(organization_bin_or_iin) = 12), -- Organization BIN/IIN
    workplace VARCHAR(255),                   -- Workplace
    invited_person_status VARCHAR(50)         -- Status of the invited person in the case
);

CREATE TABLE investigation_details (
    id SERIAL PRIMARY KEY,
    planned_investigative_actions TEXT NOT NULL, -- Planned investigative actions
    scheduled_date_time TIMESTAMP NOT NULL,      -- Scheduled date and time of investigation
    location VARCHAR(255),                       -- Investigation location
    type_of_investigation VARCHAR(50),           -- Type of investigation
    expected_outcome TEXT                        -- Expected outcome of the investigation
);
CREATE TABLE organizer_details (
    id SERIAL PRIMARY KEY,
    investigator VARCHAR(255) NOT NULL  -- Investigator information
);

CREATE TABLE business_details (
    id SERIAL PRIMARY KEY,
    is_business_related BOOLEAN NOT NULL,                  -- Is business-related
    pension_bin_or_iin VARCHAR(12) CHECK (LENGTH(pension_bin_or_iin) = 12), -- Pension BIN/IIN
    pension_workplace VARCHAR(255),                        -- Pension workplace
    entrepreneur_participation TEXT                        -- Reason and necessity for entrepreneur participation
);
CREATE TABLE defender_details (
    id SERIAL PRIMARY KEY,
    defender_iin VARCHAR(12) CHECK (LENGTH(defender_iin) = 12), -- Defender IIN, 12 digits
    defender_full_name VARCHAR(255)                             -- Defender full name
);

CREATE TABLE cards (
    id SERIAL PRIMARY KEY,
    creation_date timestamp(0) with time zone NOT NULL DEFAULT NOW(),         -- Дата создания документа
    region VARCHAR(50) NOT NULL,              -- Регион
    case_details_id INT REFERENCES case_details(id),
    person_details_id INT REFERENCES person_details(id),
    investigation_details_id INT REFERENCES investigation_details(id),
    organizer_details_id INT REFERENCES organizer_details(id),
    business_details_id INT REFERENCES business_details(id),
    defender_details_id INT REFERENCES defender_details(id)
);
