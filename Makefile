include .env

build:
	@go build -o ${BIN_PATH}/${APP_NAME} ./main.go

sqlc-generate:
	sqlc generate -f ${SQLC_CONFIG}

sqlc-clean:
	rm ${SQLC_OUT_PATH}/*.sql.go
	rm ${SQLC_OUT_PATH}/db.go
	rm ${SQLC_OUT_PATH}/querier.go
	rm ${SQLC_OUT_PATH}/model.go

docker-new-psql-container:
	@docker volume create ${APP_NAME}-volume
	@docker run --name ${APP_NAME}-postgres \
	-p ${POSTGRES_PORT}:5432 \
	-v ${APP_NAME}-volume:/var/lib/postgresql/data \
	-e POSTGRES_USER=${POSTGRES_USERNAME} \
	-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
	-d \
	${POSTGRES_IMAGE_TAG}

docker-create-db:
	docker exec ${APP_NAME}-postgres psql -U ${POSTGRES_USERNAME} -c "CREATE DATABASE ${POSTGRES_DB_NAME};"

docker-flush-db:
	docker container rm ${APP_NAME}-postgres
	docker volume rm ${APP_NAME}-volume

docker-down-db:
	docker stop ${APP_NAME}-postgres

docker-up-db:
	docker start ${APP_NAME}-postgres

migrate-create:
	@echo "Name of .sql?: "; \
    read FILENAME; \
	migrate create -ext sql -dir ${MIGRATION_PATH} -seq $${FILENAME} 

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

gen-jwt-secret:
	@openssl rand -base64 ${JWT_SECRET_LEN} > ${JWT_SECRET_OUT_PATH}

run: docker-up-db
	rm ./${APP_NAME}
	go build -o ./${APP_NAME} main.go
	chmod +x ./${APP_NAME}
	./${APP_NAME} -v v1 -c ./config/config.json -s development -h localhost -p 8000

about: ## Display info related to the build
	@echo "- Protoc version  : $(shell protoc --version)"
	@echo "- Go version      : $(shell go version)"
	@echo "- migrate version : ${shell migrate -version)}"
