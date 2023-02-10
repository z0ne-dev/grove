-- create table config
create table config
(
    id      xid   not null
        constraint config_pk
            primary key,
    user_id xid
        constraint config_user
            unique,
    value   jsonb not null
);

