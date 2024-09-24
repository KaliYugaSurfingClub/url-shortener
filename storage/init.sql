DROP TABLE user;
DROP TABLE link;
DROP TABLE click;

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
    custom_name VARCHAR(255) NOT NULL,
    clicks_count INTEGER DEFAULT 0 NOT NULL,
    last_access_time TIMESTAMP,

    expiration_date TIMESTAMP,
    clicks_to_expiration INTEGER,
    archived BOOLEAN DEFAULT FALSE NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (created_by) REFERENCES user ON DELETE SET NULL,
    UNIQUE (custom_name, created_by)
);

CREATE TABLE IF NOT EXISTS click(
    id INTEGER PRIMARY KEY,
    link_id INTEGER,

    access_time TIMESTAMP NOT NULL,
    ip VARCHAR(45) NOT NULL,
    ad_status TINYINT NOT NULL CHECK (ad_status IN (0, 1, 2)), ---started, closed, completed

    FOREIGN KEY (link_id) REFERENCES link(id) ON DELETE SET NULL
);
