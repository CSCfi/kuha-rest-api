CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    client_token TEXT NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (client_token) REFERENCES clients(client_token) ON DELETE CASCADE
);
