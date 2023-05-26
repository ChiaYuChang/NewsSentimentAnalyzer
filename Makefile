include .env

POSTGRESQL_URL="postgres://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB_NAME}?sslmode=disable"

build:
	@go build -o ${BIN_PATH} ./main.go

run:
	@go run ./main.go

sqlc-generate:
	sqlc generate -f ./config/sqlc.yml

sqlc-clean:
	rm ./internal/models/*.sql.go
	rm ./internal/models/db.go
	rm ./internal/models/models.go

docker-new-psql-container:
	@docker volume create nsa-volume
	@docker run --name nsa-postgres \
	-p ${POSTGRES_PORT}:5432 \
	-v nsa-volume:/var/lib/postgresql/data \
	-e POSTGRES_USER=${POSTGRES_USERNAME} \
	-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
	-d \
	postgres:14.6

docker-create-db:
	docker exec nsa-postgres psql -U postgres -c "CREATE DATABASE ${POSTGRES_DB_NAME};"

docker-flush-db:
	docker container rm nsa-postgres
	docker volume rm nsa-volume

docker-down-db:
	docker stop nsa-postgres

docker-up-db:
	docker start nsa-postgres

migrate-up:
	migrate -path ${MIGRATION_PATH}/ -database ${POSTGRESQL_URL} -verbose up 1

migrate-down:
	migrate -path ${MIGRATION_PATH}/ -database ${POSTGRESQL_URL} -verbose down 1

db-dump:
	pg_dump  ${POSTGRESQL_URL} -f ${SQL_SCHEME_PATH}/schema.sql --schema-only

about: ## Display info related to the build
	@echo "- Protoc version  : $(shell protoc --version)"
	@echo "- Go version      : $(shell go version)"
	@echo "- migrate version : ${shell migrate -version}"
