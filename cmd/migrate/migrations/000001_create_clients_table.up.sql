CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    client_name TEXT NOT NULL,           
    client_token TEXT UNIQUE NOT NULL,   
    role TEXT NOT NULL,                 
    created_at TIMESTAMP DEFAULT now()
);
