CREATE TABLE entrepreneur_participation (
    id SERIAL PRIMARY KEY,
    card_id INT REFERENCES cards(id) ON DELETE CASCADE,
    participation_description TEXT,
    version integer NOT NULL DEFAULT 1

);
