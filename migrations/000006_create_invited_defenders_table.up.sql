CREATE TABLE defenders (
    id SERIAL PRIMARY KEY,
    card_id INT REFERENCES cards(id) ON DELETE CASCADE,
    iin CHAR(12) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    version integer NOT NULL DEFAULT 1
);
