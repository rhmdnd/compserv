#!/bin/bash
# A common set of bash function that are useful for setting up and tearing down
# test infrastructure.

function set_container_runtime {
    if [ -z "$RUNTIME" ]; then
        if command -v podman 1>/dev/null 2>&1; then
            RUNTIME=podman
        else
            RUNTIME=docker
        fi
    fi
}

function set_database_environment_variables {
    DB_USER=${DB_USER:-"dbadmin"}
    DB_PASSWORD=${DB_PASSWORD:-"secret"}
    DB_NAME=${DB_NAME:-"compliance"}
    DB_HOST=${DB_HOST:-"localhost"}
}

function cleanup_database_container {
    $RUNTIME stop postgres
    $RUNTIME rm postgres
}

function create_database_container {
    $RUNTIME run -d --name postgres \
        -e POSTGRES_USER="$DB_USER" \
        -e POSTGRES_DB="$DB_NAME" \
        -e POSTGRES_PASSWORD="$DB_PASSWORD" \
        -p 5432:5432 \
        --health-cmd pg_isready \
        --health-interval 10s \
        --health-timeout 5s \
        --health-retries 5 \
        postgres:latest
}

function wait_for_db_init {
    health_status=""

    for _ in $(seq 1 30); do
        if [ $RUNTIME == "podman" ]; then
            healthcheck_str="{{if .State.Healthcheck}}{{print .State.Healthcheck.Status}}{{end}}"
        else
            healthcheck_str="{{if .Config.Healthcheck}}{{print .State.Health.Status}}{{end}}"
        fi

        health_status=$($RUNTIME inspect --format="$healthcheck_str" postgres)
        if [ "$health_status" == "healthy" ] ; then
            break
        fi
        sleep 3
    done

    if [ "$health_status" != "healthy" ] ; then
        echo "Failed to wait for pgsql container to come up"
        exit 1
    fi
}
