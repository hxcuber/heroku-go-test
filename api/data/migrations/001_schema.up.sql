-- noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE IF NOT EXISTS users (
    user_id BIGSERIAL PRIMARY KEY,
    user_email TEXT UNIQUE NOT NULL
);

CREATE TYPE STATUS AS ENUM ('r_blocked_s', 'r_subscribed_s', 'friends');

CREATE TABLE IF NOT EXISTS relationships (
    receiver_id BIGINT, sender_id BIGINT,
    status STATUS NOT NULL,
    PRIMARY KEY ( receiver_id, sender_id ),
    FOREIGN KEY ( sender_id ) REFERENCES users ( user_id ),
    FOREIGN KEY ( receiver_id ) REFERENCES users ( user_id )
);

