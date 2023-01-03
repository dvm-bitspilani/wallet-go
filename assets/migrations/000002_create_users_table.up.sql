CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    created TIMESTAMPTZ NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL
);
