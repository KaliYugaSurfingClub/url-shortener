CREATE TABLE person(
    id BIGSERIAL PRIMARY KEY,

    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE link(
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT REFERENCES person(id) ON DELETE SET NULL,

    original TEXT NOT NULL,
    alias VARCHAR(255) NOT NULL UNIQUE,
    custom_name VARCHAR(255) NOT NULL, --todo can be null and unique for one user
    clicks_count INTEGER DEFAULT 0 NOT NULL,
    last_access_time TIMESTAMP,

    expiration_date TIMESTAMP,
    clicks_to_expiration INTEGER,
    archived BOOLEAN DEFAULT FALSE NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    UNIQUE (custom_name, created_by)
);

CREATE TABLE click(
    id BIGSERIAL PRIMARY KEY,
    link_id BIGINT REFERENCES link(id) ON DELETE SET NULL,

    access_time TIMESTAMP NOT NULL,
    ip INET NOT NULL,
    ad_status SMALLINT NOT NULL CHECK (ad_status IN (0, 1, 2))
);
