CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    username      varchar(255) not null unique,
    password_hash varchar(255) not null
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