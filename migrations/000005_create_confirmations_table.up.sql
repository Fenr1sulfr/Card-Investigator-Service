CREATE TABLE IF NOT EXISTS card_confirmations (
    id SERIAL PRIMARY KEY,
    card_id INT NOT NULL,
    user_id INT NOT NULL,
    FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
