BEGIN;

CREATE TABLE person(
    id BIGSERIAL PRIMARY KEY,

    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,

    balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE link(
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT REFERENCES person(id) NOT NULL,

    original TEXT NOT NULL,
    alias VARCHAR(255) NOT NULL UNIQUE,
    custom_name VARCHAR(255) NOT NULL,
    clicks_count INTEGER DEFAULT 0 NOT NULL,
    last_access_time TIMESTAMP,

    expiration_date TIMESTAMP,
    clicks_to_expire INTEGER,
    archived BOOLEAN DEFAULT FALSE NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    UNIQUE (custom_name, created_by)
);

CREATE TABLE ad_video(
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT REFERENCES person(id) NOT NULL,
    uuid UUID NOT NULL,
    views_count INTEGER
);

CREATE TYPE ad_session_status_enum AS ENUM ('opened', 'closed', 'completed');

CREATE TABLE ad_session(
    id BIGSERIAL PRIMARY KEY,
    link_id BIGINT REFERENCES link(id) NOT NULL,
    ad_video_id BIGINT REFERENCES ad_video(id),

    user_agent TEXT NOT NULL,
    ip INET NULL,
    access_time TIMESTAMP NOT NULL,

    status ad_session_status_enum DEFAULT 'opened' NOT NULL
);

CREATE TABLE payment(
    id BIGSERIAL PRIMARY KEY,
    person_id BIGINT REFERENCES person(id) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE payment ADD CONSTRAINT no_direct_inserts CHECK (false);

CREATE TABLE referral_payment(
    referred_user_id BIGINT REFERENCES person(id) NOT NULL
) INHERITS (payment);

CREATE TABLE ad_payment(
    ad_session_id BIGINT REFERENCES ad_session(id) NOT NULL
) INHERITS (payment);

CREATE FUNCTION update_ad_session_status() RETURNS TRIGGER AS $$
DECLARE
    current_status ad_session_status_enum;
BEGIN
    PERFORM 1 FROM ad_session WHERE id = NEW.session_id FOR UPDATE;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Ad session with id % does not exist.', NEW.ad_session_id;
    END IF;

    SELECT status INTO current_status
    FROM ad_session
    WHERE id = NEW.ad_session_id;

    IF current_status = 'completed' THEN
        RAISE EXCEPTION 'Ad session with id % is already completed.', NEW.ad_session_id;
    END IF;

    UPDATE ad_session
    SET status = 'completed'
    WHERE id = NEW.ad_session_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_click_ad_status_on_insert_to_ad_payment
AFTER INSERT ON ad_payment
FOR EACH ROW EXECUTE FUNCTION update_ad_session_status();

CREATE FUNCTION update_link_and_ad_video_on_session_open() RETURNS TRIGGER AS $$
BEGIN
    PERFORM 1 FROM link WHERE id = NEW.link_id FOR UPDATE;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Link with id % does not exist', NEW.link_id;
    END IF;

    PERFORM 1 FROM ad_video WHERE id = NEW.ad_video_id FOR UPDATE;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Ad video with id % does not exist', NEW.ad_video_id;
    END IF;

    UPDATE link
    SET clicks_count = clicks_count + 1, last_access_time = NEW.access_time
    WHERE id = NEW.link_id;

    UPDATE ad_video
    SET views_count = views_count + 1
    WHERE id = NEW.ad_video_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_link_and_ad_video_on_session_open
AFTER INSERT ON ad_session
FOR EACH ROW EXECUTE FUNCTION update_link_and_ad_video_on_session_open();

COMMIT;
