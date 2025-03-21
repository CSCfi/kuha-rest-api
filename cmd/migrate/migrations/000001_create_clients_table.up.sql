CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    client_name TEXT NOT NULL,           -- Human-readable name (e.g., UTV, FIS)
    client_token TEXT UNIQUE NOT NULL,   -- Hashed token 
    role TEXT NOT NULL,                  -- e.g. 'admin'
    created_at TIMESTAMP DEFAULT now()
);
