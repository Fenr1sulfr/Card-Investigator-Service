CREATE TABLE IF NOT EXISTS card_confirmations (
    id SERIAL PRIMARY KEY,
    card_id INT REFERENCES cards(id),
    user_id INT REFERENCES users(id),
    UNIQUE (card_id, user_id) -- предотвращает дублирование
);
