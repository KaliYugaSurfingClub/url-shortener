CREATE TABLE IF NOT EXISTS user(
    id INTEGER PRIMARY KEY,

    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS link(
    id INTEGER PRIMARY KEY,
    created_by INTEGER,

    original TEXT NOT NULL,
    alias VARCHAR(255) NOT NULL UNIQUE,
    clicks_count INTEGER DEFAULT 0 NOT NULL,
    last_access_time TIMESTAMP,

    expiration_date TIMESTAMP,
    max_clicks INTEGER,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (created_by) REFERENCES user
);

CREATE TABLE IF NOT EXISTS click(
    id INTEGER PRIMARY KEY,
    link_id INTEGER,

    access_time TIMESTAMP NOT NULL,
    ip VARCHAR(45) NOT NULL,
    full_ad BOOLEAN NOT NULL,

    FOREIGN KEY (link_id) REFERENCES link(id)
);
