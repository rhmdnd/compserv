#!/bin/bash

MIGRATION_DIR=$(cd $(dirname "${BASH_SOURCE:-$0}") && pwd)
source $MIGRATION_DIR/functions

set_container_runtime
set_database_envs
create_database_container
wait_for_db_init
