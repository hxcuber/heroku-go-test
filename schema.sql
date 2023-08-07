-- noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE IF NOT EXISTS users (
    user_email TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS friends (
    user_1_email TEXT NOT NULL,
    user_2_email TEXT NOT NULL,
    PRIMARY KEY(user_1_email, user_2_email)
);

CREATE TABLE IF NOT EXISTS update_subscribed (
     user_email TEXT NOT NULL,
     subscribed_to_email TEXT NOT NULL,
     PRIMARY KEY (user_email, subscribed_to_email)
);

CREATE TABLE IF NOT EXISTS update_blocked (
    user_email TEXT NOT NULL,
    blocked_by_email TEXT NOT NULL,
    PRIMARY KEY (user_email, blocked_by_email)
);
