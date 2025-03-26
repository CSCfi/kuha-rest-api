-- clients
CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY,
    client_name TEXT NOT NULL,
    client_token TEXT UNIQUE NOT NULL,
    role TEXT[] NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

-- refresh_tokens
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    client_token TEXT NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (client_token) REFERENCES clients(client_token) ON DELETE CASCADE
);

-- revoked_tokens
CREATE TABLE IF NOT EXISTS revoked_tokens (
    client_token TEXT PRIMARY KEY,
    revoked_at TIMESTAMP DEFAULT now()
);

-- revoked_refresh_tokens
CREATE TABLE IF NOT EXISTS revoked_refresh_tokens (
    token TEXT PRIMARY KEY,
    revoked_at TIMESTAMP DEFAULT now()
);

-- token_logs
CREATE TABLE IF NOT EXISTS token_logs (
    id SERIAL PRIMARY KEY,
    client_token TEXT NOT NULL REFERENCES clients(client_token) ON DELETE CASCADE,
    token_type TEXT NOT NULL,
    action TEXT NOT NULL,
    token TEXT,
    ip_address TEXT,
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT now()
);