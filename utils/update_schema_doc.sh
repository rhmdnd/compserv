#!/bin/bash
# This script updates the database schema documentation. This script should be
# invoked each time the schema is changed.

UTILS_DIR=$(cd "$(dirname "${BASH_SOURCE:-$0}")" && pwd)
TOP_DIR=$(dirname "$UTILS_DIR")
# shellcheck source=utils/functions.sh
source "${UTILS_DIR}/functions.sh"

set_container_runtime
set_database_environment_variables
trap cleanup_database_container EXIT
create_database_container
wait_for_db_init


POSTGRESQL_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST/$DB_NAME?sslmode=disable"
"$TOP_DIR/tools/migrate" -database "$POSTGRESQL_URL" -path "$TOP_DIR/migrations" up
podman exec -it postgres pg_dump --username dbadmin --schema-only compliance > "${TOP_DIR}/migrations/schema.sql"
