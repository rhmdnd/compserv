#!/bin/bash
# This script is a test runner for testing the compliance service against a
# real database.
#
# Stand up the database and run the tests in the same script so that we can
# cleanup the database regardless of the test outcome. Otherwise, if we were to
# create the database in the Makefile we would need to handle that gracefully
# there so we can cleanup the container if the tests fail. For now, let's do it
# all from this script and just invoke it from the Makefile in an effort to
# keep the Makefile cleaner.


UTILS_DIR=$(cd "$(dirname "${BASH_SOURCE:-$0}")" && pwd)
TOP_DIR=$(dirname "$UTILS_DIR")
# shellcheck source=utils/functions.sh
source "${UTILS_DIR}/functions.sh"

set_container_runtime
set_database_environment_variables
trap cleanup_database_container EXIT
create_database_container
wait_for_db_init

go test -v "${TOP_DIR}/tests/..."
