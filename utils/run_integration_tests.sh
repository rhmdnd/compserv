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

# Make sure we run integration tests serially by isolating them to a single
# process. This isn't ideal for all tests, but since we're testing against a
# real databse that needs to be upgraded and downgraded, we could introduce
# transient issues if multiple tests are attempting to modify state at the same
# time. This constraint is specific to the database tests. When we increase
# integration coverage to include more tests (e.g., that exercise the API
# layer) then we should consider organizing the integration tests into separate
# modules. In the interest of performance, we should keep the tests that run
# serially as small as possible.
go test -v -p 1 "${TOP_DIR}/tests/..."
