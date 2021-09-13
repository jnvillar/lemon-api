CONTAINER_NAME ?= lemon_mysql
DB_PASSWORD = your_secret
DB_PORT = 3306
DB_NAME = lemon
DB ?= 'mysql://root:$(DB_PASSWORD)@tcp(localhost:$(DB_PORT))/lemon?query'
MIGRATE := PATH="$$PWD/bin:$$PATH" \
           bin/migrate -source="file:$$PWD/storage/mysql/migrations"


build:
	go build -o bin/lemon main.go

run: build
	./bin/lemon serve --sql_password=$(DB_PASSWORD)

test:
	go test ./...

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

create-mysql-container: delete-mysql-container
	docker run \
	  -d -e MYSQL_ROOT_PASSWORD=$(DB_PASSWORD) -e  MYSQL_ROOT_HOST="172.17.0.1" -p '$(DB_PORT):$(DB_PORT)' \
	  --name $(CONTAINER_NAME) mysql/mysql-server:latest

delete-mysql-container:
	docker rm -f $(CONTAINER_NAME)


database/create:
	docker exec $(CONTAINER_NAME) mysql -uroot -p$(DB_PASSWORD) -e "create database $(DB_NAME)";

database/delete:
	docker exec $(CONTAINER_NAME) mysql -uroot -p$(DB_PASSWORD) -e "drop database $(DB_NAME)";

database/fill:
	docker exec -i $(CONTAINER_NAME) mysql -uroot -p$(DB_PASSWORD) lemon < "./storage/mysql/dump/dump.sql";

migrate/build:
	go build -tags 'mysql' -o bin/migrate \
		github.com/golang-migrate/migrate/v4/cmd/migrate

database/migrations/up: migrate/build
	$(MIGRATE) -database=$(DB) up


database/migrations/down: migrate/build
	$(MIGRATE) -database=$(DB) down