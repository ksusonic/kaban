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

create table if not exists tasks
(
    id          serial primary key,
    board_id    int          not null references boards (id) on delete cascade,
    author_id   int          not null,
    assignee_id int          not null,
    title       varchar(255) not null,
    description text         not null,
    priority    smallint     not null,
    due_date    timestamp,
    status      int          not null default 0,
    created_at  timestamp    not null default now(),
    updated_at  timestamp    not null default now(),
    deleted_at  timestamp             default null
);

create index if not exists idx_tasks_board_id on tasks (board_id);
create index if not exists idx_tasks_deleted_at on tasks (deleted_at);

