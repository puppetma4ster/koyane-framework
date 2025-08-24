#!/bin/bash

DB="wordists.db"
SQLSETUP="$(dirname "$0")/tableSetup/setup.sql" # path to database tables configuration

SQLENTITIES_DIR="$(dirname "$0")/SQLentities" # folder for db entries

PRINT="$(dirname "$0")/printAll.sql"  # print whole DB an the and

# loading Setup & Data
sqlite3 "$DB" < "$SQLSETUP"

# loading every .sql in SQLentities folder
for file in "$SQLENTITIES_DIR"/*.sql; do
    echo "Importing $file..."
    sqlite3 "$DB" < "$file"
done

sqlite3 "$DB" < "$PRINT"