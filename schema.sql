CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    merchant TEXT NOT NULL,
    card_or_pass TEXT NOT NULL,
    amount TEXT NOT NULL,
    date TEXT NOT NULL
);
