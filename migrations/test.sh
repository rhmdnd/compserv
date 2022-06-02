#!/bin/bash

MIGRATION_DIR=$(cd $(dirname "${BASH_SOURCE:-$0}") && pwd)
source $MIGRATION_DIR/functions

if [ -z $MIGRATE ]; then
    MIGRATE=migrate
fi

set_container_runtime
set_database_envs

trap cleanup_database_container EXIT

create_database_container
wait_for_db_init

POSTGRESQL_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:5432/$DB_NAME?sslmode=disable"

$MIGRATE -database $POSTGRESQL_URL -path migrations up
echo "y" | $MIGRATE -database $POSTGRESQL_URL -path migrations down
