/*TODO create users table */

CREATE TABLE test (
    id SERIAL PRIMARY KEY,
    status BOOLEAN,
    score INT,
    email TEXT NOT NULL
);

INSERT INTO test (status, score, email) VALUES
(true, 85, 'alice@example.com'),
(false, 72, 'bob@example.com'),
(true, 90, 'carol@example.com'),
(false, 60, 'dave@example.com'),
(true, 88, 'eve@example.com');