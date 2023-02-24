gooseConnectionString = "postgres://postgres:postgres@localhost:5432/psn_discount?sslmode=disable"
gooseTable = "__goose_db_version"

add_migration:
	goose -dir ./migrations/postgres create $(name) sql

migrate:
	goose -dir ./migrations/postgres -table $(gooseTable) postgres $(gooseConnectionString) up

run:
	go run main.go -c config.yaml

build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -a -installsuffix cgo -o binfile

dockerbuild: build
	docker build -t psn_discount_bot .

dockerrun:
	docker run -d --name psn_discount psn_discount_bot
