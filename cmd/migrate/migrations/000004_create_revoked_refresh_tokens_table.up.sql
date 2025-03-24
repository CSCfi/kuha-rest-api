CREATE TABLE revoked_refresh_tokens (
    id SERIAL PRIMARY KEY,
    token TEXT UNIQUE NOT NULL,                 -- The refresh token (stored in plaintext)
    revoked_at TIMESTAMP DEFAULT now()          -- When it was revoked
);
