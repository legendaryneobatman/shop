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

create table brands
(
    id               serial primary key,
    name             varchar(255)        not null,
    slug             varchar(255) unique not null,
    description      text,
    logo_url         varchar(500),
    website_url      varchar(500),

    is_active        boolean   default true,
    is_featured      boolean   default false,

    display_order    int       default 0,

    meta_title       varchar(255),
    meta_description varchar(500),
    meta_keywords    varchar(500),

    created_at       timestamp default now(),
    updated_at       timestamp default now()
);


CREATE TABLE categories
(
    id               serial primary key,
    parent_id        int references categories (id) on delete cascade,
    name             varchar(255)        not null,
    slug             varchar(255) unique not null,
    description      text,
    image_url        varchar(500),
    display_order    int       default 0,
    is_active        bool      default true,
    meta_title       varchar(255),
    meta_description varchar(500),
    meta_keywords    varchar(500),
    created_at       timestamp default now(),
    updated_at       timestamp default now(),
    level            int       default 0,
    path             varchar(1000)
);

create table products
(
    id                  serial primary key,
    sku                 varchar(50) unique  not null,
    name                varchar(255)        not null,
    slug                varchar(255) unique not null,
    description         text,
    short_description   varchar(500),
    price               decimal(10, 2)      not null,
    old_price           decimal(10, 2),

    category_id         int                 references categories (id) on delete set null,
    brand_id            int                 references brands (id) on delete set null,

    quantity            int                 not null default 0,
    low_stock_threshold int                          default 5,

    main_image_url      varchar(500),
    weight              decimal(8, 2),
    dimensions          varchar(100),

    is_active           boolean                      default true,
    is_featured         boolean                      default false,
    is_new              boolean                      default false,
    rating              decimal(3, 2)                default 0.0,
    review_count        int                          default 0,

    created_at          timestamp                    default now(),
    updated_at          timestamp                    default now(),
    meta_title          varchar(255),
    meta_description    varchar(500),
    meta_keywords       varchar(500)
);