# Sample Makefile for a Go project

APP_NAME = OneApp
GO_VERSION = 1.25-bookworm
# GO_VERSION = alpine3.22
PWD_PATH = $(CURDIR)

.PHONY: up

t1:
	echo ${PWD_PATH}

test:
	docker run --rm -it --name=test \
	 -v ${PWD_PATH}/data:/data \
	 -w /data \
	 datacatering/duckdb:v1.3.2

duck:
	duckdb ${PWD_PATH}/data/db/data.db -ui


tidy:
	docker run --rm -it --name=test \
	 -v ${PWD_PATH}/server:/go/src/app \
	 -w /go/src/app \
	 --pid=host \
	 golang:$(GO_VERSION) \
	 go mod tidy


server-build:
	cd server && \
	docker build -f Dockerfile.debug -t server-app-debug .

# 	docker run -d -p 8080:8080 -p 4000:4000 --name server-app-debug server-app-debug
# 	docker run -d -it --rm -p 8080:8080 -p 4000:4000 --name server-app-debug server-app-debug sleep infinity

up:
	docker compose  -f docker-compose.yml up -d

up-build: server-build
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build --watch

down:
	docker compose down -v


#  make goose CMD="status"    
#  https://github.com/pressly/goose
goose:
	cd data && goose $(CMD)


# all: build

# build:
# 	go build -o $(APP_NAME) .

# run: build
# 	./$(APP_NAME)

# test:
# 	go test ./...

# clean:
# 	rm -f $(APP_NAME)
