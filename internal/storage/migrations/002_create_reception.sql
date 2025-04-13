CREATE TABLE IF NOT EXISTS receptions (
    id UUID NOT NULL PRIMARY KEY,
    pvz_id UUID NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    date_time TIMESTAMP NOT NULL DEFAULT NOW(),
    status TEXT NOT NULL CHECK (status IN ('in_progress', 'close'))
);

