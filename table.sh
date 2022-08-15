#!/bin/sh
  DB_USER='steve'
  DB_PASS='secret'
  DB='vault'
cat << EOF | psql --username "$DB_USER" --password "$DB_PASS" --dbname "$DB"
CREATE TABLE vault (id serial PRIMARY KEY, secret varchar, duration varchar, uuid varchar);

EOF
