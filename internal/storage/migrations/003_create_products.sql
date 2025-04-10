CREATE TABLE IF NOT EXISTS products (
    id UUID NOT NULL PRIMARY KEY,
    date_time TIMESTAMP NOT NULL DEFAULT NOW(),
    type TEXT NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    reception_id UUID NOT NULL REFERENCES receptions(id) ON DELETE CASCADE
);

