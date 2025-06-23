#!/usr/bin/env bash
set -e

psql -v ON_ERROR_STOP=0 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER usr_crud WITH LOGIN PASSWORD 'Passw0rd';
	CREATE DATABASE db_crud WITH OWNER usr_crud;

EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "db_crud" <<-EOSQL
    create table books (
        id serial not null unique,
        title varchar(255) not null unique,
        author varchar(255) not null,
        publish_date timestamp not null default now(),
        rating int not null
    );

    ALTER TABLE public.books OWNER TO usr_crud;

EOSQL