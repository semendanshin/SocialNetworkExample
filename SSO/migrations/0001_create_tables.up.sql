CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id              UUID PRIMARY KEY,                             -- Primary key with UUID type
    username        VARCHAR(255) NOT NULL UNIQUE,                 -- Username with a unique constraint
    email           VARCHAR(255) NOT NULL UNIQUE,                 -- Email with a unique constraint
    hashed_password VARCHAR(255) NOT NULL,                        -- Password column
    created_at      TIMESTAMP    NOT NULL DEFAULT NOW(),          -- Timestamp for record creation
    updated_at      TIMESTAMP    NOT NULL DEFAULT NOW()           -- Timestamp for record update
);

CREATE INDEX idx_users_username ON users(username);

CREATE INDEX idx_users_email ON users(email);
