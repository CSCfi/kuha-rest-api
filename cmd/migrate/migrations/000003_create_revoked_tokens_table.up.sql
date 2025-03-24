CREATE TABLE revoked_tokens (
    id SERIAL PRIMARY KEY,
    client_token TEXT UNIQUE NOT NULL,
    revoked_at TIMESTAMP DEFAULT now()
);
