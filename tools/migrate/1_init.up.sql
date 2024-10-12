BEGIN;

CREATE TABLE person
(
    id            BIGSERIAL PRIMARY KEY,

    username      VARCHAR(50)                         NOT NULL UNIQUE,
    email         VARCHAR(255)                        NOT NULL UNIQUE,
    password_hash VARCHAR(255)                        NOT NULL,

    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE balance
(
    id        BIGSERIAL PRIMARY KEY,

    person_id BIGINT REFERENCES person (id) NOT NULL,
    amount    DECIMAL(10, 2)                NOT NULL DEFAULT 0.00 CHECK (amount >= 0)
);

CREATE TABLE link
(
    id               BIGSERIAL PRIMARY KEY,
    created_by       BIGINT REFERENCES person (id)       NOT NULL,

    original         TEXT                                NOT NULL,
    alias            VARCHAR(255)                        NOT NULL UNIQUE,
    custom_name      VARCHAR(255)                        NOT NULL,

    expiration_date  TIMESTAMP,
    clicks_to_expire INTEGER,
    archived         BOOLEAN   DEFAULT FALSE             NOT NULL,

    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    UNIQUE (custom_name, created_by)
);

CREATE TABLE ad_video
(
    id              BIGSERIAL PRIMARY KEY,
    created_by      BIGINT REFERENCES person (id)       NOT NULL,

    title           VARCHAR(255)                        NOT NULL,
    expiration_date TIMESTAMP,

    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    UNIQUE (title, created_by)
);

CREATE TABLE ad_session
(
    id          BIGSERIAL PRIMARY KEY,
    link_id     BIGINT REFERENCES link (id) NOT NULL,
    ad_video_id BIGINT REFERENCES ad_video (id),

    user_agent  TEXT                        NOT NULL,
    ip          INET                        NULL,
    access_time TIMESTAMP                   NOT NULL
);

CREATE TABLE payment
(
    id         BIGSERIAL PRIMARY KEY,
    person_id  BIGINT REFERENCES person (id) NOT NULL,
    amount     DECIMAL(10, 2)                NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE referral_payment
(
    id               BIGSERIAL PRIMARY KEY,
    payment_id       BIGINT REFERENCES payment (id) NOT NULL,
    referred_user_id BIGINT REFERENCES person (id)  NOT NULL
);

CREATE TABLE ad_payment
(
    id            BIGSERIAL PRIMARY KEY,
    payment_id    BIGINT REFERENCES payment (id)    NOT NULL,
    ad_session_id BIGINT REFERENCES ad_session (id) NOT NULL
);

COMMIT;
