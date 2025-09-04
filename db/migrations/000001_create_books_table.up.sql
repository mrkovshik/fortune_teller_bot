create table if not exists books
(
    id     bigserial
        constraint books_pk
            primary key,
    title  varchar not null,
    author varchar not null,
    format varchar not null,
    lang   varchar not null
);