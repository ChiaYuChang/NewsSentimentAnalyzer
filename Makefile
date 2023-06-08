include .env

STATE=testing
APP_REPOSITORY="github.com/ChiaYuChang/NewsSentimentAnalyzer"
POSTGRESQL_URL="postgres://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB_NAME}_${STATE}?sslmode=disable"

build:
	@go build -o ${BIN_PATH} ./main.go

run:
	@go run ./main.go

sqlc-generate:
	sqlc generate -f ./config/sqlc.yml

sqlc-clean:
	rm ./internal/server/model/*.sql.go
	rm ./internal/server/model/db.go
	rm ./internal/server/model/querier.go
	rm ./internal/server/model/model.go

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
	docker exec nsa-postgres psql -U postgres -c "CREATE DATABASE ${POSTGRES_DB_NAME}_${STATE};"

docker-flush-db:
	docker container rm nsa-postgres
	docker volume rm nsa-volume

docker-down-db:
	docker stop nsa-postgres

docker-up-db:
	docker start nsa-postgres

migrate-up:
	migrate -path ${MIGRATION_PATH}/ -database ${POSTGRESQL_URL} -verbose up 1

migrate-up-all:
	migrate -path ${MIGRATION_PATH}/ -database ${POSTGRESQL_URL} -verbose up

migrate-down:
	migrate -path ${MIGRATION_PATH}/ -database ${POSTGRESQL_URL} -verbose down 1

migrate-down-all:
	migrate -path ${MIGRATION_PATH}/ -database ${POSTGRESQL_URL} -verbose down

migrate-drop:
	migrate -path ${MIGRATION_PATH}/ -database ${POSTGRESQL_URL} -verbose goto 1
	migrate -path ${MIGRATION_PATH}/ -database ${POSTGRESQL_URL} -verbose drop

db-dump:
	pg_dump  ${POSTGRESQL_URL} -f ${SQL_SCHEME_PATH}/schema.sql --schema-only

mockgen-store:
	@mockgen -destination internal/server/model/mockdb/store.go \
	${APP_REPOSITORY}/internal/server/model Store

mockgen-tokenmaker:
	@mockgen -destination pkgs/tokenMaker/mockTokenMaker/tokenMaker.go \
	${APP_REPOSITORY}/pkgs/tokenMaker TokenMaker,Payload

about: ## Display info related to the build
	@echo "- Protoc version  : $(shell protoc --version)"
	@echo "- Go version      : $(shell go version)"
	@echo "- migrate version : ${shell migrate -version}"
