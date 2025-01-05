create table if not exists users
(
    id         serial primary key,
    username   varchar(255) not null,
    first_name varchar(255),
    last_name  varchar(255),
    avatar_url text,
    created_at timestamp    not null default now(),
    updated_at timestamp    not null default now()
);

create table if not exists telegram_users
(
    telegram_id bigint primary key,
    user_id     integer not null references users (id) on delete cascade
);

create index if not exists idx_users_user_id on telegram_users (user_id);
