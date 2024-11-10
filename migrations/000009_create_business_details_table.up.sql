CREATE TABLE business_details (
    id SERIAL PRIMARY KEY,
    is_business_related BOOLEAN NOT NULL,                  -- Относится ли к бизнесу (справочник);
    pension_bin_or_iin VARCHAR(12) CHECK (LENGTH(pension_bin_or_iin) = 12), -- Пенсионные BIN/IIN
    pension_workplace VARCHAR(255),                        -- Место работы (пенсионные отчисления) (автоподтягивание последнего места работы с базы пенсионных отчислений по ИИН вызываемого, отображается пользователю с ролью Аналитик СД);
    entrepreneur_participation TEXT                        -- Обоснование и необходимость участия предпринимателя (обязательное поле, ручной ввод);
);