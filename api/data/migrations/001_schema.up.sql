-- noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE IF NOT EXISTS users (
    user_id BIGSERIAL PRIMARY KEY,
    user_email TEXT UNIQUE NOT NULL
);

CREATE TYPE SUBSCRIPTION_STATUS AS ENUM ('r_blocked_s', 'r_subscribed_s', 'none');

CREATE TABLE IF NOT EXISTS relationships (
    sender_id BIGINT,
    receiver_id BIGINT,
    friends BOOL NOT NULL,
    status SUBSCRIPTION_STATUS NOT NULL,
    PRIMARY KEY(sender_id, receiver_id),
    FOREIGN KEY(sender_id) REFERENCES users(user_id),
    FOREIGN KEY(receiver_id) REFERENCES users(user_id)
);
