CREATE TABLE investigation_details (
    id SERIAL PRIMARY KEY,
    planned_investigative_actions TEXT NOT NULL, -- Планируемые следственные действия (обязательное поле, ручнойввод);
    scheduled_date_time TIMESTAMP NOT NULL,      -- Дата и время проведения (календарный и временной выбор);
    location VARCHAR(255),                       -- Место проведения (справочник);
    type_of_investigation VARCHAR(50),           -- Виды планируемого следствия (справочник);
    expected_outcome TEXT                        -- Результат от планируемого следственного действия (обязательное поле, ручной ввод).
    
);