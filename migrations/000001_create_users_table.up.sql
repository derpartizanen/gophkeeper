CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    user_id  uuid DEFAULT gen_random_uuid () primary key,
    username varchar(128) not null unique,
    security_key text not null
);
