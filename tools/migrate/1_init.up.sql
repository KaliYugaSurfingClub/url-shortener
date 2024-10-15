BEGIN;

CREATE TYPE ad_type AS ENUM ('video', 'file');
CREATE TYPE click_status AS ENUM ('completed', 'opened');

CREATE TABLE person
(
    id            BIGSERIAL PRIMARY KEY,

    username      VARCHAR(50)                         NOT NULL UNIQUE,
    email         VARCHAR(255)                        NOT NULL UNIQUE,
    password_hash VARCHAR(255)                        NOT NULL,

    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE ad_source
(
    id              BIGSERIAL PRIMARY KEY,
    person_id       BIGINT REFERENCES person (id)       NOT NULL,

    title           VARCHAR(255)                        NOT NULL,
    type            ad_type                             NOT NULL,
    expiration_date TIMESTAMP                           NOT NULL,
    start_date      TIMESTAMP                           NOT NULL,

    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    UNIQUE (title, person_id)
);

CREATE TABLE link
(
    id          BIGSERIAL PRIMARY KEY,
    person_id   BIGINT REFERENCES person (id)       NOT NULL,

    original    TEXT                                NOT NULL,
    alias       VARCHAR(255)                        NOT NULL UNIQUE,
    custom_name VARCHAR(255)                        NOT NULL,
    archived    BOOLEAN   DEFAULT FALSE             NOT NULL,

    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    UNIQUE (custom_name, person_id)
);

CREATE TABLE click
(
    id           BIGSERIAL PRIMARY KEY,
    link_id      BIGINT REFERENCES link (id)      ON DELETE SET NULL,
    ad_source_id BIGINT REFERENCES ad_source (id) NOT NULL,

    user_agent   TEXT                             NOT NULL,
    ip           INET                             NOT NULL,
    access_time  TIMESTAMP                        NOT NULL
);

CREATE TABLE balance_transaction
(
    id         BIGSERIAL PRIMARY KEY,
    person_id  BIGINT REFERENCES person (id)       NOT NULL,
    amount     DECIMAL(10, 2)                      NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE referral_reward
(
    id                     BIGSERIAL PRIMARY KEY,
    balance_transaction_id BIGINT REFERENCES balance_transaction (id) NOT NULL UNIQUE
    ---todo
);

CREATE TABLE click_reward
(
    id                     BIGSERIAL PRIMARY KEY,
    balance_transaction_id BIGINT REFERENCES balance_transaction (id) NOT NULL UNIQUE,
    click_id               BIGINT REFERENCES click (id)               NOT NULL UNIQUE
);

CREATE TABLE ad_source_purchase
(
    id                     BIGSERIAL PRIMARY KEY,
    balance_transaction_id BIGINT REFERENCES balance_transaction (id) NOT NULL UNIQUE,
    ad_source_id           BIGSERIAL REFERENCES ad_source (id)        NOT NULL UNIQUE
);

CREATE INDEX index_click_link_id ON click (link_id);

COMMIT;
