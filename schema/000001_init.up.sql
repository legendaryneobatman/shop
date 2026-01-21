CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    username      varchar(255) not null unique,
    password_hash varchar(255) not null,

    email         varchar(255),
    avatar_url    varchar(255),
    phone         varchar(255),
    role          int, -- енум по интам усл. -1 admin, 1 user, 2 moderator etc.
    is_active     bool      default true,
    created_at    timestamp default CURRENT_TIMESTAMP,
    updated_at    timestamp default CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END
$$ language plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE refresh_tokens
(
    id         serial                                      not null unique,
    user_id    int references users (id) on delete cascade not null,
    token_hash varchar(255)                                not null,
    expires_at timestamp                                   not null,
    ip_address varchar(255),
    user_agent varchar,
    revoked    bool default false
);

CREATE TABLE lists
(
    id          serial                                      not null unique,
    title       varchar(255)                                not null,
    description varchar(255),

    user_id     int references users (id) on delete cascade not null
);

CREATE TABLE todos
(
    id          serial       not null unique,
    title       varchar(255) not null,
    description varchar(255),
    done        boolean      not null default false,

    list_id     int references lists (id) on delete cascade
);