CREATE TABLE organizer_details (
    id SERIAL PRIMARY KEY,
    investigator VARCHAR(255) NOT NULL  -- Следователь (автоподтягивание с личного кабинета);
);