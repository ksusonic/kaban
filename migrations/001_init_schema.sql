create table if not exists users
(
    id          serial primary key,
    username    varchar(255)  not null,
    telegram_id bigint unique not null,
    first_name  varchar(255),
    last_name   varchar(255),
    avatar_url  text,
    created_at  timestamp default NOW(),
    updated_at  timestamp default NOW()
);

CREATE INDEX if not exists idx_users_telegram_id ON users (telegram_id);
