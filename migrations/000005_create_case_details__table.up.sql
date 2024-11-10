CREATE TABLE case_details (
    id SERIAL PRIMARY KEY,
    case_number VARCHAR(15) NOT NULL CHECK (LENGTH(case_number) = 15), -- Номер дела, 15-digit format
    registration_date TIMESTAMP NOT NULL,          -- Дата регистрации дела
    criminal_code_article VARCHAR(50),             -- Статья УК
    case_decision TEXT,                            -- Решение по делу
    case_summary TEXT,                             -- "Краткая фабула(чтобы это не значило)"
    relation_to_event TEXT                         -- Отношение вызываемого к событию и субьекту
    
);