CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,
    card_id INT REFERENCES cards(id) ON DELETE CASCADE,
    bin_or_iin CHAR(12) NOT NULL,
    workplace VARCHAR(100),
    version integer NOT NULL DEFAULT 1    
);
