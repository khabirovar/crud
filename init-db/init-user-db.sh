#!/usr/bin/env bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER usr_crud WITH LOGIN PASSWORD 'Passw0rd';
	CREATE DATABASE db_crud WITH OWNER usr_crud;
    

EOSQL