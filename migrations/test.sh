#!/bin/bash

MIGRATIONS_DIR=$(cd "$(dirname "${BASH_SOURCE:-$0}")" && pwd)
TOP_DIR=$(dirname "$MIGRATIONS_DIR")
# shellcheck source=utils/functions.sh
source "${TOP_DIR}/utils/functions.sh"

if [ -z "$MIGRATE" ]; then
    MIGRATE=migrate
fi

set_container_runtime
set_database_environment_variables
trap cleanup_database_container EXIT
create_database_container
wait_for_db_init

POSTGRESQL_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:5432/$DB_NAME?sslmode=disable"

$MIGRATE -database "$POSTGRESQL_URL" -path migrations up
echo "y" | $MIGRATE -database "$POSTGRESQL_URL" -path migrations down
