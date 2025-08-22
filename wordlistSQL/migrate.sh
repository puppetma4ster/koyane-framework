#!/bin/bash

DB="wl_database.db"

SQLSETUP="$(dirname "$0")/setup.sql"
SQLDATA="$(dirname "$0")/data.sql"

sqlite3 "$DB" < "$SQLSETUP"
sqlite3 "$DB" < "$SQLDATA"
