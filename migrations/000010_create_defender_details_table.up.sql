CREATE TABLE defender_details (
    id SERIAL PRIMARY KEY,
    defender_iin VARCHAR(12) CHECK (LENGTH(defender_iin) = 12), -- ИИН защитника (ручной ввод, ФЛК 12 цифр);
    defender_full_name VARCHAR(255)                             -- ФИО защитника (автоподтягивание по ИИН защитника);
);