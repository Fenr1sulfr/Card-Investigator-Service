CREATE TABLE person_details (
    id SERIAL PRIMARY KEY,
    invited_person_iin VARCHAR(12) NOT NULL CHECK (LENGTH(invited_person_iin) = 12), -- ИИН, 12 digits
    invited_person_full_name VARCHAR(255),    -- ФИО приглашенного
    invited_person_position VARCHAR(100),     -- Позиция
    organization_bin_or_iin VARCHAR(12) NOT NULL CHECK (LENGTH(organization_bin_or_iin) = 12), -- Организации BIN/IIN
    workplace VARCHAR(255),                   -- Место работы
    invited_person_status VARCHAR(50),       -- Статус приглашенного
    version integer NOT NULL DEFAULT 1
);