create extension if not exists postgis;
create extension if not exists "uuid-ossp";

create table if not exists app_users
(
    id uuid default public.uuid_generate_v4() not null
        constraint pk_app_users
            primary key,
    first_name text,
    last_name text,
    email text,
    password text not null
);

create unique index if not exists app_users_email_uindex on app_users (email);
