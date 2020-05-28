begin;

create table "user"
(
    id bigserial    not null constraint user_pk primary key,
    name       varchar(255) not null,
    email      varchar(255) not null,
    status     varchar(255) default 'disabled'::character varying,
    active_key text,
    created_at timestamp,
    updated_at timestamp
);

create unique index user_email_uindex on "user" (email);

create table files
(
    id bigserial not null constraint files_pk primary key,
    original_name varchar(255) not null,
    external_url  varchar(255) not null,
    created_at    timestamp    default now(),
    updated_at    timestamp,
    user_id       bigint not null,
    size          bigint       default 0,
    type          varchar(255) default 'file'::character varying not null
);


END;
