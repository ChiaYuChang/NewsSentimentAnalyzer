#! bin/bash

function migrate_up() {
    if [[ ! -z "$1" ]]; then
        local POSTGRES_PASSWORD="$(cat $1)"
        migrate -path ${MIGRATION_PATH}/ -database postgres://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${DB_NAME}?sslmode=${SSL_MODE} up
    else
        echo "[Warning] passing password without encryption is not recommended"
        echo "[Warning] set DB_PASSWORD_FILE variable"
        migrate -path ${MIGRATION_PATH}/ -database postgres://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${DB_NAME}?sslmode=${SSL_MODE} up
    fi
}

migrate_up ${POSTGRES_PASSWORD_FILE}
