CREATE TABLE IF NOT EXISTS files (
    id serial PRIMARY KEY,
    parent_id INTEGER REFERENCES files(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    size INTEGER NOT NULL,
    type TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
