BEGIN;

DROP TABLE IF EXISTS click;
DROP TABLE IF EXISTS link;
DROP TABLE IF EXISTS person;

CREATE TABLE person(
    id BIGSERIAL PRIMARY KEY,

    username VARCHAR(50) NOT NULL UNIQUE CHECK (username <> ''),
    email VARCHAR(255) NOT NULL UNIQUE CHECK (email <> ''),
    password_hash VARCHAR(255) NOT NULL CHECK (password_hash <> ''),

    balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE link(
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT REFERENCES person(id) NOT NULL,

    original TEXT NOT NULL CHECK (original <> ''),
    alias VARCHAR(255) NOT NULL UNIQUE CHECK (alias <> ''),
    custom_name VARCHAR(255) NOT NULL CHECK (custom_name <> ''),
    clicks_count INTEGER DEFAULT 0 NOT NULL CHECK (clicks_count >= 0),
    last_access_time TIMESTAMP,

    expiration_date TIMESTAMP CHECK (expiration_date IS NULL OR expiration_date > CURRENT_TIMESTAMP),
    clicks_to_expire INTEGER CHECK (clicks_to_expire IS NULL OR clicks_to_expire > 0),
    archived BOOLEAN DEFAULT FALSE NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    UNIQUE (custom_name, created_by)
);

CREATE TABLE click(
    id BIGSERIAL PRIMARY KEY,
    link_id BIGINT REFERENCES link(id) NOT NULL,

    user_agent TEXT NOT NULL,
    ip INET NULL,
    access_time TIMESTAMP NOT NULL,

    ad_status SMALLINT NOT NULL CHECK (ad_status IN (0, 1, 2))
);

COMMIT;
