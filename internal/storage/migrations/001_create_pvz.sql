CREATE TABLE IF NOT EXISTS pvz (
    id UUID NOT NULL PRIMARY KEY,
    registration_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    city TEXT NOT NULL 
    CONSTRAINT chk_pvz_city CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань'))
);