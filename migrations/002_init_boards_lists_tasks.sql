create table if not exists boards
(
    id         serial primary key,
    name       varchar(255) not null,
    slug       varchar(255) not null unique,
    owner_id   int          not null, -- nullable to allow user deletion
    created_at timestamp    not null default now(),
    updated_at timestamp    not null default now(),
    deleted_at timestamp             default null
);

create index if not exists idx_boards_owner_id on boards (owner_id);
create index if not exists idx_boards_deleted_at on boards (deleted_at);

create table if not exists board_members
(
    board_id     int       not null references boards (id) on delete cascade,
    user_id      int       not null references users (id) on delete cascade,
    access_level int       not null,
    added_at     timestamp not null default now(),
    updated_at   timestamp not null default now(),
    deleted_at   timestamp          default null
);

create unique index if not exists idx_board_members_board_id_user_id on board_members (board_id, user_id);
create index if not exists idx_board_members_deleted_at on board_members (deleted_at);

create table if not exists lists
(
    id         serial primary key,
    board_id   int          not null references boards (id) on delete cascade,
    name       varchar(255) not null,
    position   int          not null unique,
    created_at timestamp    not null default now(),
    updated_at timestamp    not null default now()
);

create index if not exists idx_lists_board_id on lists (board_id);

create table if not exists tasks
(
    id          serial primary key,
    list_id     int          not null references lists (id) on delete cascade,
    title       varchar(255) not null,
    metadata    jsonb,
    description text,
    position    int          not null unique,
    due_date    timestamp,
    status      int          not null default 0,
    created_at  timestamp    not null default now(),
    updated_at  timestamp    not null default now(),
    deleted_at  timestamp             default null
);

create index if not exists idx_tasks_list_id on tasks (list_id);
create index if not exists idx_tasks_deleted_at on tasks (deleted_at);

