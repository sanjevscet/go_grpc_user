CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO users (name, email, username, password, role)
VALUES
('Sanjeev', 'sanjeev@example.com', 'sanjeev', 'sanjeev123', 'admin'),
('Rahul', 'rahul@example.com', 'rahul', 'rahul123', 'user'),
('Aman', 'aman@example.com', 'aman', 'aman123', 'user'),
('Harshit', 'harshit@example.com', 'harshit', 'harshit123', 'user')
ON CONFLICT (email) DO NOTHING;