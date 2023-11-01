CREATE TABLE IF NOT EXISTS files (
    id serial PRIMARY KEY,
    parent_id INTEGER REFERENCES files(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    size INTEGER NOT NULL,
    type TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
