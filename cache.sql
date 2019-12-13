CREATE TABLE cache (
    id SERIAL PRIMARY KEY,
    query text NOT NULL,
    result JSON,
    created_at timestamp with time zone NOT NULL
);