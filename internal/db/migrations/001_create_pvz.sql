CREATE TABLE IF NOT EXISTS pvz (
    id UUID NOT NULL PRIMARY KEY,
    registration_date TIMESTAMP NOT NULL DEFAULT NOW(),
    city TEXT NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань'))
);