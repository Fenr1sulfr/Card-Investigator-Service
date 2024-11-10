CREATE TABLE basic_info (
    id SERIAL PRIMARY KEY,
    registration_number VARCHAR(50) UNIQUE NOT NULL, -- Unique registration number
    creation_date TIMESTAMP NOT NULL,                -- Document creation date
    region VARCHAR(50) NOT NULL                      -- Region information
    
);