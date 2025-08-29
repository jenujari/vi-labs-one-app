# Sample Makefile for a Go project

APP_NAME = OneApp
GO_VERSION = 1.25-bookworm
# GO_VERSION = alpine3.22
PWD_PATH = $(CURDIR)

# Load .env file
ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: up

t1:
	@echo "Running $(POSTGRES_DB)"

# test:
# 	cd headless && docker build -t headless-test . && docker run --rm -it -p 7878:7878 headless-test

tidy-server:
	docker run --rm -it --name=tidy-server \
	 -v ${PWD_PATH}/server:/go/src/app \
	 -w /go/src/app \
	 --pid=host \
	 golang:$(GO_VERSION) \
	 go mod tidy

tidy-headless:
	docker run --rm -it --name=tidy-headless \
	 -v ${PWD_PATH}/headless:/go/src/app \
	 -w /go/src/app \
	 --pid=host \
	 golang:$(GO_VERSION) \
	 go mod tidy

tidy: tidy-headless tidy-server

server-build:
	cd server && \
	docker build -f Dockerfile.debug -t server-app-debug .

# 	docker run -d -p 8080:8080 -p 4000:4000 --name server-app-debug server-app-debug
# 	docker run -d -it --rm -p 8080:8080 -p 4000:4000 --name server-app-debug server-app-debug sleep infinity

up:
	docker compose -f docker-compose.yml up -d

up-infra:
	docker compose  -f docker-compose.yml up -d timescaledb pgadmin

# server-build
up-build:  up-infra
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build --watch

down:
	docker compose down -v

#  make goose CMD="status"    
#  https://github.com/pressly/goose
goose:
	cd data && goose $(CMD)

sqlc:
	docker run --rm -v $(PWD_PATH)/server/sqlc:/src -w /src sqlc/sqlc generate

zip-pgadmin:
	7z a -tzip -p$(SEVEN_ZIP_PASSWORD) -mem=AES256 pgadmin.zip ./data/pg-admin

zip-pgsql:
	7z a -tzip -p$(SEVEN_ZIP_PASSWORD) -mem=AES256 pgsql.zip ./data/pgsql

uzip-pgadmin:
	7z x pgadmin.zip -p$(SEVEN_ZIP_PASSWORD) -o./data/pg-admin/restoration

uzip-pgsql:
	7z x pgsql.zip -p$(SEVEN_ZIP_PASSWORD) -o./data/pgsql/restoration

headless-test:
	cd headless && go test ./...

server-test:
	cd server && go test ./...

go-test: server-test headless-test

