include .env

build:
	@go build -o ${BIN_PATH}/${APP_NAME} ./main.go && \
		chmod +x ${BIN_PATH}/${APP_NAME}

build-wasm:
	@GOARCH=wasm GOOS=js go build -o ${WASM_PATH}/${WASM_NAME} ./thirdparty/wasm/main.go

sqlc-generate:
	sqlc generate -f ${SQLC_CONFIG}

sqlc-clean:
	rm ${SQLC_OUT_PATH}/*.sql.go
	rm ${SQLC_OUT_PATH}/db.go
	rm ${SQLC_OUT_PATH}/querier.go
	rm ${SQLC_OUT_PATH}/model.go

docker-new-psql-container:
	@docker volume create ${APP_NAME}-postgres-volume
	@docker run --name ${APP_NAME}-postgres \
	-p ${POSTGRES_PORT}:5432 \
	-v ${APP_NAME}-postgres-volume:/var/lib/postgresql/data \
	-e POSTGRES_USER=${POSTGRES_USERNAME} \
	-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
	-d \
	${POSTGRES_IMAGE_TAG}

docker-new-redis-container:
	@docker volume create ${APP_NAME}-redis-volume
	@docker run --name ${APP_NAME}-redis \
	-p ${REDIS_PORT}:6379 \
	-v ${APP_NAME}-redis-volume:/data \
	-e REDIS_ARGS="--requirepass ${REDIS_PASSWORD}" \
	-d \
	${REDIS_IMAGE_TAG}

docker-create-db:
	docker exec ${APP_NAME}-postgres psql -U ${POSTGRES_USERNAME} -c "CREATE DATABASE ${POSTGRES_DB_NAME};"

docker-flush-db:
	docker container rm ${APP_NAME}-postgres
	docker container rm ${APP_NAME}-redis
	docker volume rm ${APP_NAME}-postgres-volume
	docker volume rm ${APP_NAME}-redis-volume

docker-down-db:
	@docker stop ${APP_NAME}-postgres
	@docker stop ${APP_NAME}-redis

docker-up-db:
	@docker start ${APP_NAME}-postgres
	@docker start ${APP_NAME}-redis

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
	@mockgen -destination internal/model/mockdb/store.go \
	${APP_REPOSITORY}/internal/model Store

mockgen-tokenmaker:
	@mockgen -destination pkgs/tokenMaker/mockTokenMaker/tokenMaker.go \
	${APP_REPOSITORY}/pkgs/tokenMaker TokenMaker,Payload

mockgen-parser:
	@mockgen -destination internal/server/parser/mockParser/parser.go \
	${APP_REPOSITORY}/internal/server/parser Parser,StdParseProcess

gen-jwt-secret:
	@openssl rand -base64 ${JWT_SECRET_LEN} > ${JWT_SECRET_OUT_PATH}

gen-private-key:
	@openssl ecparam -genkey \
	-name secp384r1 \
	-out ${KEY_PATH}/${PRIVATE_KEY_NAME}

gen-public-key: gen-private-key
	@openssl req -new -x509 -sha256 \
	-key ${KEY_PATH}/${PRIVATE_KEY_NAME} \
	-out ${KEY_PATH}/${PUBLIC_KEY_NAME} \
	-config ssl.conf \
	-days 365

# gen-ssl-key:
# 	openssl req -x509 -new -nodes -sha256 -utf8 \
# 	-days 3650 \
# 	-newkey rsa:2048 \
# 	-keyout ${KEY_PATH}/${PRIVATE_KEY_NAME} \
# 	-out ${KEY_PATH}/${PUBLIC_KEY_NAME} \
# 	-config ssl.conf

build-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	${PROTO_SRC_DIR}/*.proto \

clean-proto:
	rm -rf ./proto/*pb.go

build-lang-detector:
	@go build -o bin/languageDetectorServer/languageDetectorServer ./cmd/languageDetectorServer/main.go && \
		chmod +x bin/languageDetectorServer/languageDetectorServer

build-milvus-health-check:
	@go build -o bin/milvusHealthCheck/milvusHealthCheck ./cmd/milvusHealthCheck/main.go && \
		chmod +x bin/milvusHealthCheck/milvusHealthCheck

build-news-sentiment-analyzer:
	@go build -o ./${APP_NAME} ./main.go && \
		chmod +x ./${APP_NAME}

run: docker-up-db build build-wasm build-news-sentiment-analyzer
	./${APP_NAME} -v v1 -c ./config/config.json -s development -h localhost -p 8001

build: build-lang-detector build-milvus-health-check build-news-sentiment-analyzer

start-lang-detect-server:
	@./bin/languageDetectorServer/languageDetectorServer -o ./bin/languageDetectorServer/log.json

clean:
	@rm ./${APP_NAME}
	@rm bin/milvusHealthCheck/milvusHealthCheck
	@rm bin/languageDetectorServer/languageDetectorServer

about: ## Display info related to the build
	@echo "- Protoc version  : $(shell protoc --version)"
	@echo "- Go version      : $(shell go version)"
	@echo "- migrate version : ${shell migrate -version)}"
