#!/bin/bash

MIGRATION_DIR=$(cd $(dirname "${BASH_SOURCE:-$0}") && pwd)
source $MIGRATION_DIR/functions

set_container_runtime
cleanup_database_container
