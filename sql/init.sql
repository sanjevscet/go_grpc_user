CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);

INSERT INTO users (name, email) VALUES
    ('Rahul', 'rahul@example.com'),
    ('Aman', 'aman@example.com'),
    ('Sanjeev', 'sanjeev@example.com'),
    ('Harshit', 'harshit@example.com')
ON CONFLICT (email) DO NOTHING;
